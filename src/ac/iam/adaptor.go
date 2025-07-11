/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package iam

import (
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

// NotEnoughLayer TODO
var NotEnoughLayer = fmt.Errorf("not enough layer")

// AdaptAuthOptions TODO
func AdaptAuthOptions(a *meta.ResourceAttribute) (ActionID, []types.Resource, error) {

	var action ActionID

	action, err := ConvertResourceAction(a.Type, a.Action, a.BusinessID)
	if err != nil {
		return "", nil, err
	}

	// convert different cmdb resource's to resource's registered to iam
	rscType, err := ConvertResourceType(a.Type, a.BusinessID)
	if err != nil {
		return "", nil, err
	}

	resource, err := GenIamResource(action, *rscType, a)
	if err != nil {
		return "", nil, err
	}

	return action, resource, nil
}

var ccIamResTypeMap = map[meta.ResourceType]TypeID{
	meta.Business:                 Business,
	meta.BizSet:                   BizSet,
	meta.Project:                  Project,
	meta.Model:                    SysModel,
	meta.ModelUnique:              SysModel,
	meta.ModelAttributeGroup:      SysModel,
	meta.ModelModule:              BizTopology,
	meta.ModelSet:                 BizTopology,
	meta.MainlineInstance:         BizTopology,
	meta.MainlineInstanceTopology: BizTopology,
	meta.MainlineModel:            TypeID(""),
	meta.ModelTopology:            TypeID(""),
	meta.ModelClassification:      SysModelGroup,
	meta.AssociationType:          SysAssociationType,
	meta.ModelAssociation:         SysModel,
	meta.MainlineModelTopology:    TypeID(""),
	meta.ModelInstanceTopology:    SkipType,
	meta.CloudAreaInstance:        SysCloudArea,
	meta.HostInstance:             Host,
	meta.HostFavorite:             SkipType,
	meta.Process:                  BizProcessServiceInstance,
	meta.DynamicGrouping:          BizCustomQuery,
	meta.AuditLog:                 SysAuditLog,
	meta.SystemBase:               TypeID(""),
	meta.UserCustom:               UserCustom,
	meta.ProcessServiceTemplate:   BizProcessServiceTemplate,
	meta.ProcessServiceCategory:   BizProcessServiceCategory,
	meta.ProcessServiceInstance:   BizProcessServiceInstance,
	meta.BizTopology:              BizTopology,
	meta.SetTemplate:              BizSetTemplate,
	meta.OperationStatistic:       SysOperationStatistic,
	meta.HostApply:                BizHostApply,
	meta.ResourcePoolDirectory:    SysResourcePoolDirectory,
	meta.CloudAccount:             SysCloudAccount,
	meta.CloudResourceTask:        SysCloudResourceTask,
	meta.EventWatch:               SysEventWatch,
	meta.ConfigAdmin:              TypeID(""),
	meta.SystemConfig:             TypeID(""),
	meta.KubeCluster:              TypeID(""),
	meta.KubeNode:                 TypeID(""),
	meta.KubeNamespace:            TypeID(""),
	meta.KubeWorkload:             TypeID(""),
	meta.KubeDeployment:           TypeID(""),
	meta.KubeStatefulSet:          TypeID(""),
	meta.KubeDaemonSet:            TypeID(""),
	meta.KubeGameStatefulSet:      TypeID(""),
	meta.KubeGameDeployment:       TypeID(""),
	meta.KubeCronJob:              TypeID(""),
	meta.KubeJob:                  TypeID(""),
	meta.KubePodWorkload:          TypeID(""),
	meta.KubePod:                  TypeID(""),
	meta.KubeContainer:            TypeID(""),
	meta.FieldTemplate:            FieldGroupingTemplate,
	meta.FulltextSearch:           TypeID(""),
	meta.IDRuleIncrID:             TypeID(""),
	meta.FullSyncCond:             TypeID(""),
	meta.GeneralCache:             GeneralCache,
}

// ConvertResourceType convert resource type from CMDB to IAM
func ConvertResourceType(resourceType meta.ResourceType, businessID int64) (*TypeID, error) {
	var iamResourceType TypeID

	switch resourceType {
	case meta.ModelAttribute:
		if businessID > 0 {
			iamResourceType = BizCustomField
		} else {
			iamResourceType = SysModel
		}
		return &iamResourceType, nil
	case meta.NetDataCollector:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	iamResourceType, exists := ccIamResTypeMap[resourceType]
	if exists {
		return &iamResourceType, nil
	}

	if IsCMDBSysInstance(resourceType) {
		iamResourceType = TypeID(resourceType)
		return &iamResourceType, nil
	}

	return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
}

// ConvertResourceAction convert resource action from CMDB to IAM
func ConvertResourceAction(resourceType meta.ResourceType, action meta.Action, businessID int64) (ActionID, error) {
	if action == meta.SkipAction {
		return Skip, nil
	}

	convertAction := action
	switch action {
	case meta.CreateMany:
		convertAction = meta.Create
	case meta.FindMany:
		convertAction = meta.Find
	case meta.DeleteMany:
		convertAction = meta.Delete
	case meta.UpdateMany:
		convertAction = meta.Update
	}

	switch resourceType {
	case meta.ModelAttribute, meta.ModelAttributeGroup:
		switch convertAction {
		case meta.Delete, meta.Update, meta.Create:
			if businessID > 0 {
				return EditBusinessCustomField, nil
			} else {
				return EditSysModel, nil
			}
		}
	case meta.HostInstance:
		switch convertAction {
		case meta.Update:
			if businessID > 0 {
				return EditBusinessHost, nil
			} else {
				return EditResourcePoolHost, nil
			}
		case meta.Find:
			if businessID > 0 {
				return ViewBusinessResource, nil
			} else {
				return ViewResourcePoolHost, nil
			}
		}
	}

	if IsCMDBSysInstance(resourceType) {
		return ConvertSysInstanceActionID(resourceType, convertAction)
	}

	if _, exist := resourceActionMap[resourceType]; exist {
		actionID, ok := resourceActionMap[resourceType][convertAction]
		if ok && actionID != Unsupported {
			return actionID, nil
		}
	}

	return Unsupported, fmt.Errorf("unsupported type %s action: %s", resourceType, action)
}

// ConvertSysInstanceActionID convert system instances action from CMDB to IAM
func ConvertSysInstanceActionID(resourceType meta.ResourceType, action meta.Action) (ActionID, error) {
	var actionType ActionType
	switch action {
	case meta.Create:
		actionType = Create
	case meta.Update:
		actionType = Edit
	case meta.Delete:
		actionType = Delete
	case meta.Find:
		actionType = View
	default:
		return Unsupported, fmt.Errorf("unsupported action: %s", action)
	}
	id := strings.TrimPrefix(string(resourceType), meta.CMDBSysInstTypePrefix)
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return Unsupported, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	return ActionID(fmt.Sprintf("%s_%s%s", actionType, IAMSysInstTypePrefix, id)), nil
}

var resourceActionMap = map[meta.ResourceType]map[meta.Action]ActionID{
	meta.ModelAttributeGroup: {
		meta.Delete:   EditSysModel,
		meta.Update:   EditSysModel,
		meta.Create:   EditSysModel,
		meta.Find:     ViewSysModel,
		meta.FindMany: ViewSysModel,
	},
	meta.ModelUnique: {
		meta.Delete:   EditSysModel,
		meta.Update:   EditSysModel,
		meta.Create:   EditSysModel,
		meta.Find:     ViewSysModel,
		meta.FindMany: ViewSysModel,
	},
	meta.Business: {
		meta.Archive:              ArchiveBusiness,
		meta.Create:               CreateBusiness,
		meta.Update:               EditBusiness,
		meta.Find:                 FindBusiness,
		meta.ViewBusinessResource: ViewBusinessResource,
	},
	meta.BizSet: {
		meta.Create:       CreateBizSet,
		meta.Update:       EditBizSet,
		meta.Delete:       DeleteBizSet,
		meta.Find:         ViewBizSet,
		meta.AccessBizSet: AccessBizSet,
	},
	meta.DynamicGrouping: {
		meta.Delete:   DeleteBusinessCustomQuery,
		meta.Update:   EditBusinessCustomQuery,
		meta.Create:   CreateBusinessCustomQuery,
		meta.Find:     ViewBusinessResource,
		meta.FindMany: ViewBusinessResource,
		meta.Execute:  ViewBusinessResource,
	},
	meta.MainlineModel: {
		meta.Find:   Skip,
		meta.Create: EditBusinessLayer,
		meta.Delete: EditBusinessLayer,
	},
	meta.ModelTopology: {
		meta.Find:              EditModelTopologyView,
		meta.Update:            EditModelTopologyView,
		meta.ModelTopologyView: ViewModelTopo,
	},
	meta.MainlineModelTopology: {
		meta.Find: Skip,
	},
	meta.Process: {
		meta.Find:   Skip,
		meta.Create: EditBusinessServiceInstance,
		meta.Delete: EditBusinessServiceInstance,
		meta.Update: EditBusinessServiceInstance,
	},
	meta.HostInstance: {
		meta.MoveResPoolHostToBizIdleModule: ResourcePoolHostTransferToBusiness,
		meta.MoveResPoolHostToDirectory:     ResourcePoolHostTransferToDirectory,
		meta.MoveBizHostFromModuleToResPool: BusinessHostTransferToResourcePool,
		meta.AddHostToResourcePool:          CreateResourcePoolHost,
		meta.Create:                         CreateResourcePoolHost,
		meta.Delete:                         DeleteResourcePoolHost,
		meta.MoveHostToAnotherBizModule:     HostTransferAcrossBusiness,
		meta.Find:                           ViewResourcePoolHost,
		meta.FindMany:                       ViewResourcePoolHost,
		meta.ManageHostAgentID:              ManageHostAgentID,
	},
	meta.ProcessServiceCategory: {
		meta.Delete: DeleteBusinessServiceCategory,
		meta.Update: EditBusinessServiceCategory,
		meta.Create: CreateBusinessServiceCategory,
		meta.Find:   Skip,
	},
	meta.ProcessServiceInstance: {
		meta.Delete:   DeleteBusinessServiceInstance,
		meta.Update:   EditBusinessServiceInstance,
		meta.Create:   CreateBusinessServiceInstance,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ProcessServiceTemplate: {
		meta.Delete:   DeleteBusinessServiceTemplate,
		meta.Update:   EditBusinessServiceTemplate,
		meta.Create:   CreateBusinessServiceTemplate,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.SetTemplate: {
		meta.Delete:   DeleteBusinessSetTemplate,
		meta.Update:   EditBusinessSetTemplate,
		meta.Create:   CreateBusinessSetTemplate,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelModule: {
		meta.Delete:   DeleteBusinessTopology,
		meta.Update:   EditBusinessTopology,
		meta.Create:   CreateBusinessTopology,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelSet: {
		meta.Delete:   DeleteBusinessTopology,
		meta.Update:   EditBusinessTopology,
		meta.Create:   CreateBusinessTopology,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.MainlineInstance: {
		meta.Delete:   DeleteBusinessTopology,
		meta.Update:   EditBusinessTopology,
		meta.Create:   CreateBusinessTopology,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.MainlineInstanceTopology: {
		meta.Delete: Skip,
		meta.Update: Skip,
		meta.Create: Skip,
		meta.Find:   Skip,
	},
	meta.HostApply: {
		meta.Create:           EditBusinessHostApply,
		meta.Update:           EditBusinessHostApply,
		meta.Delete:           EditBusinessHostApply,
		meta.Find:             Skip,
		meta.DefaultHostApply: ViewBusinessResource,
	},
	meta.ResourcePoolDirectory: {
		meta.Delete:                DeleteResourcePoolDirectory,
		meta.Update:                EditResourcePoolDirectory,
		meta.Create:                CreateResourcePoolDirectory,
		meta.AddHostToResourcePool: CreateResourcePoolHost,
		meta.Find:                  Skip,
	},
	meta.CloudAreaInstance: {
		meta.Delete:   DeleteCloudArea,
		meta.Update:   EditCloudArea,
		meta.Create:   CreateCloudArea,
		meta.Find:     ViewCloudArea,
		meta.FindMany: ViewCloudArea,
	},
	meta.CloudAccount: {
		meta.Delete:   DeleteCloudAccount,
		meta.Update:   EditCloudAccount,
		meta.Create:   CreateCloudAccount,
		meta.Find:     FindCloudAccount,
		meta.FindMany: FindCloudAccount,
	},
	meta.CloudResourceTask: {
		meta.Delete: DeleteCloudResourceTask,
		meta.Update: EditCloudResourceTask,
		meta.Create: CreateCloudResourceTask,
		meta.Find:   FindCloudResourceTask,
	},
	meta.Model: {
		meta.Delete:   DeleteSysModel,
		meta.Update:   EditSysModel,
		meta.Create:   CreateSysModel,
		meta.Find:     ViewSysModel,
		meta.FindMany: ViewSysModel,
	},
	meta.AssociationType: {
		meta.Delete:   DeleteAssociationType,
		meta.Update:   EditAssociationType,
		meta.Create:   CreateAssociationType,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelClassification: {
		meta.Delete:   DeleteModelGroup,
		meta.Update:   EditModelGroup,
		meta.Create:   CreateModelGroup,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.OperationStatistic: {
		meta.Create:   EditOperationStatistic,
		meta.Delete:   EditOperationStatistic,
		meta.Update:   EditOperationStatistic,
		meta.Find:     FindOperationStatistic,
		meta.FindMany: FindOperationStatistic,
	},
	meta.AuditLog: {
		meta.Find:     FindAuditLog,
		meta.FindMany: FindAuditLog,
	},
	meta.SystemBase: {
		meta.ModelTopologyView:      EditModelTopologyView,
		meta.ModelTopologyOperation: EditBusinessLayer,
	},
	meta.EventWatch: {
		meta.WatchHost:             WatchHostEvent,
		meta.WatchHostRelation:     WatchHostRelationEvent,
		meta.WatchBiz:              WatchBizEvent,
		meta.WatchSet:              WatchSetEvent,
		meta.WatchModule:           WatchModuleEvent,
		meta.WatchProcess:          WatchProcessEvent,
		meta.WatchCommonInstance:   WatchCommonInstanceEvent,
		meta.WatchMainlineInstance: WatchMainlineInstanceEvent,
		meta.WatchInstAsst:         WatchInstAsstEvent,
		meta.WatchBizSet:           WatchBizSetEvent,
		meta.WatchPlat:             WatchPlatEvent,
		meta.WatchKubeCluster:      WatchKubeClusterEvent,
		meta.WatchKubeNode:         WatchKubeNodeEvent,
		meta.WatchKubeNamespace:    WatchKubeNamespaceEvent,
		meta.WatchKubeWorkload:     WatchKubeWorkloadEvent,
		meta.WatchKubePod:          WatchKubePodEvent,
		meta.WatchProject:          WatchProjectEvent,
	},
	meta.UserCustom: {
		meta.Find:   Skip,
		meta.Update: Skip,
		meta.Delete: Skip,
		meta.Create: Skip,
	},
	meta.ModelAssociation: {
		meta.Find:     ViewSysModel,
		meta.FindMany: ViewSysModel,
		meta.Update:   EditSysModel,
		meta.Delete:   EditSysModel,
		meta.Create:   EditSysModel,
	},
	meta.ModelInstanceTopology: {
		meta.Find:   Skip,
		meta.Update: Skip,
		meta.Delete: Skip,
		meta.Create: Skip,
	},
	meta.ModelAttribute: {
		meta.Find:   ViewSysModel,
		meta.Update: EditSysModel,
		meta.Delete: DeleteSysModel,
		meta.Create: CreateSysModel,
	},
	meta.HostFavorite: {
		meta.Find:   Skip,
		meta.Update: Skip,
		meta.Delete: Skip,
		meta.Create: Skip,
	},

	meta.ProcessTemplate: {
		meta.Find:   Skip,
		meta.Delete: DeleteBusinessServiceTemplate,
		meta.Update: EditBusinessServiceTemplate,
		meta.Create: CreateBusinessServiceTemplate,
	},
	meta.BizTopology: {
		meta.Find:   Skip,
		meta.Update: EditBusinessTopology,
		meta.Delete: DeleteBusinessTopology,
		meta.Create: CreateBusinessTopology,
	},
	// unsupported resource actions for now
	meta.NetDataCollector: {
		meta.Find:   Unsupported,
		meta.Update: Unsupported,
		meta.Delete: Unsupported,
		meta.Create: Unsupported,
	},
	meta.InstallBK: {
		meta.Update: Skip,
	},
	// TODO: confirm this
	meta.SystemConfig: {
		meta.FindMany: Skip,
		meta.Find:     Skip,
		meta.Update:   Skip,
		meta.Delete:   Skip,
		meta.Create:   Skip,
	},
	meta.ConfigAdmin: {

		// reuse GlobalSettings permissions
		meta.Find:   Skip,
		meta.Update: GlobalSettings,
		meta.Delete: Unsupported,
		meta.Create: Unsupported,
	},
	meta.KubeCluster: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerCluster,
		meta.Delete: DeleteContainerCluster,
		meta.Create: CreateContainerCluster,
	},
	meta.KubeNode: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerNode,
		meta.Delete: DeleteContainerNode,
		meta.Create: CreateContainerNode,
	},
	meta.KubeNamespace: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerNamespace,
		meta.Delete: DeleteContainerNamespace,
		meta.Create: CreateContainerNamespace,
	},
	meta.KubeWorkload: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubeDeployment: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubeStatefulSet: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubeDaemonSet: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubeGameStatefulSet: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubeGameDeployment: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubeCronJob: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubeJob: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubePodWorkload: {
		meta.Find:   ViewBusinessResource,
		meta.Update: EditContainerWorkload,
		meta.Delete: DeleteContainerWorkload,
		meta.Create: CreateContainerWorkload,
	},
	meta.KubePod: {
		meta.Find:   ViewBusinessResource,
		meta.Delete: DeleteContainerPod,
		meta.Create: CreateContainerPod,
	},
	meta.KubeContainer: {
		meta.Find: ViewBusinessResource,
	},
	meta.Project: {
		meta.Find:   ViewProject,
		meta.Update: EditProject,
		meta.Delete: DeleteProject,
		meta.Create: CreateProject,
	},
	meta.FulltextSearch: {
		meta.Find: UseFulltextSearch,
	},
	meta.FieldTemplate: {
		meta.Create: CreateFieldGroupingTemplate,
		meta.Find:   ViewFieldGroupingTemplate,
		meta.Update: EditFieldGroupingTemplate,
		meta.Delete: DeleteFieldGroupingTemplate,
	},
	meta.IDRuleIncrID: {
		meta.Update: EditIDRuleIncrID,
	},
	meta.FullSyncCond: {
		meta.Create: CreateFullSyncCond,
		meta.Find:   ViewFullSyncCond,
		meta.Update: EditFullSyncCond,
		meta.Delete: DeleteFullSyncCond,
	},
	meta.GeneralCache: {
		meta.Find: ViewGeneralCache,
	},
}

// ParseIamPathToAncestors TODO
func ParseIamPathToAncestors(iamPath []string) ([]metadata.IamResourceInstance, error) {
	instances := make([]metadata.IamResourceInstance, 0)
	for _, path := range iamPath {
		pathItemArr := strings.Split(strings.Trim(path, "/"), "/")
		for _, pathItem := range pathItemArr {
			typeAndID := strings.Split(pathItem, ",")
			if len(typeAndID) != 2 {
				return nil, fmt.Errorf("pathItem %s invalid", pathItem)
			}
			id := typeAndID[1]
			if id == "*" {
				continue
			}
			instances = append(instances, metadata.IamResourceInstance{
				Type:     typeAndID[0],
				TypeName: ResourceTypeIDMap[TypeID(typeAndID[0])],
				ID:       id,
			})
		}
	}
	return instances, nil
}

// GenIAMDynamicResTypeID 生成IAM侧资源的的dynamic resource typeID
func GenIAMDynamicResTypeID(modelID int64) TypeID {
	return TypeID(fmt.Sprintf("%s%d", IAMSysInstTypePrefix, modelID))
}

// GenCMDBDynamicResType 生成CMDB侧资源的的dynamic resourceType
func GenCMDBDynamicResType(modelID int64) meta.ResourceType {
	return meta.ResourceType(fmt.Sprintf("%s%d", meta.CMDBSysInstTypePrefix, modelID))
}

// genDynamicResourceType generate dynamic resourceType
func genDynamicResourceType(obj metadata.Object) ResourceType {
	return ResourceType{
		ID:      GenIAMDynamicResTypeID(obj.ID),
		Name:    obj.ObjectName,
		NameEn:  obj.ObjectID,
		Parents: nil,
		ProviderConfig: ResourceConfig{
			Path: "/auth/v3/find/resource",
		},
		Version: 1,
	}
}

// genDynamicResourceTypes generate dynamic resourceTypes
func genDynamicResourceTypes(objects []metadata.Object) []ResourceType {
	resourceTypes := make([]ResourceType, 0)
	for _, obj := range objects {
		resourceTypes = append(resourceTypes, genDynamicResourceType(obj))
	}
	return resourceTypes
}

// genIAMDynamicInstanceSelection generate IAM dynamic instanceSelection
func genIAMDynamicInstanceSelection(modelID int64) InstanceSelectionID {
	return InstanceSelectionID(fmt.Sprintf("%s%d", IAMSysInstTypePrefix, modelID))
}

// genDynamicInstanceSelection generate dynamic instanceSelection
func genDynamicInstanceSelection(obj metadata.Object) InstanceSelection {
	return InstanceSelection{
		ID:     genIAMDynamicInstanceSelection(obj.ID),
		Name:   obj.ObjectName,
		NameEn: obj.ObjectID,
		ResourceTypeChain: []ResourceChain{{
			SystemID: SystemIDCMDB,
			ID:       GenIAMDynamicResTypeID(obj.ID),
		}},
	}
}

// genDynamicInstanceSelections generate dynamic instanceSelections
func genDynamicInstanceSelections(objects []metadata.Object) []InstanceSelection {
	instanceSelections := make([]InstanceSelection, 0)
	for _, obj := range objects {
		instanceSelections = append(instanceSelections, genDynamicInstanceSelection(obj))
	}
	return instanceSelections
}

// genDynamicAction generate dynamic action
// Note: view action must be in the first place
func genDynamicAction(obj metadata.Object) []DynamicAction {
	return []DynamicAction{
		genDynamicViewAction(obj),
		genDynamicCreateAction(obj),
		genDynamicEditAction(obj),
		genDynamicDeleteAction(obj),
	}
}

// GenDynamicActionID generate dynamic ActionID
func GenDynamicActionID(actionType ActionType, modelID int64) ActionID {
	return ActionID(fmt.Sprintf("%s_%s%d", actionType, IAMSysInstTypePrefix, modelID))
}

// genDynamicViewAction generate dynamic view action
func genDynamicViewAction(obj metadata.Object) DynamicAction {
	return DynamicAction{
		ActionID:     GenDynamicActionID(View, obj.ID),
		ActionType:   View,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "查看"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "view", obj.ObjectID, "instance"),
	}
}

// genDynamicCreateAction generate dynamic create action
func genDynamicCreateAction(obj metadata.Object) DynamicAction {
	return DynamicAction{
		ActionID:     GenDynamicActionID(Create, obj.ID),
		ActionType:   Create,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "新建"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "create", obj.ObjectID, "instance"),
	}
}

// genDynamicEditAction generate dynamic edit action
func genDynamicEditAction(obj metadata.Object) DynamicAction {
	return DynamicAction{
		ActionID:     GenDynamicActionID(Edit, obj.ID),
		ActionType:   Edit,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "编辑"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "edit", obj.ObjectID, "instance"),
	}
}

// genDynamicDeleteAction generate dynamic delete action
func genDynamicDeleteAction(obj metadata.Object) DynamicAction {
	return DynamicAction{
		ActionID:     GenDynamicActionID(Delete, obj.ID),
		ActionType:   Delete,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "删除"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "delete", obj.ObjectID, "instance"),
	}
}

// genDynamicActionSubGroup 动态的按模型生成动作分组作为‘模型实例管理’分组的subGroup
func genDynamicActionSubGroup(obj metadata.Object) ActionGroup {
	actions := genDynamicAction(obj)
	actionWithIDs := make([]ActionWithID, len(actions))
	for idx, action := range actions {
		actionWithIDs[idx] = ActionWithID{ID: action.ActionID}
	}
	return ActionGroup{
		Name:    obj.ObjectName,
		NameEn:  obj.ObjectID,
		Actions: actionWithIDs,
	}
}

// genDynamicActionIDs generate dynamic model actionIDs
func genDynamicActionIDs(object metadata.Object) []ActionID {
	actions := genDynamicAction(object)
	actionIDs := make([]ActionID, len(actions))
	for idx, action := range actions {
		actionIDs[idx] = action.ActionID
	}
	return actionIDs
}

// genDynamicActions generate dynamic model actions
func genDynamicActions(objects []metadata.Object) []ResourceAction {
	resActions := make([]ResourceAction, 0)
	for _, obj := range objects {
		relatedResource := []RelateResourceType{
			{
				SystemID:    SystemIDCMDB,
				ID:          GenIAMDynamicResTypeID(obj.ID),
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
				SelectionMode: modeAll,
				InstanceSelections: []RelatedInstanceSelection{{
					SystemID: SystemIDCMDB,
					ID:       genIAMDynamicInstanceSelection(obj.ID),
				}},
			},
		}

		actions := genDynamicAction(obj)
		var relatedActions []ActionID
		for _, action := range actions {
			switch action.ActionType {
			case View:
				resActions = append(resActions, ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 View,
					RelatedActions:       nil,
					RelatedResourceTypes: nil,
					Version:              1,
				})
				relatedActions = []ActionID{action.ActionID}

			case Create:
				resActions = append(resActions, ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 Create,
					RelatedResourceTypes: nil,
					RelatedActions:       nil,
					Version:              1,
				})
			case Edit:
				resActions = append(resActions, ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 Edit,
					RelatedActions:       relatedActions,
					Version:              1,
					RelatedResourceTypes: relatedResource,
				})

			case Delete:
				resActions = append(resActions, ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 Delete,
					RelatedResourceTypes: relatedResource,
					RelatedActions:       relatedActions,
					Version:              1,
				})
			default:
				return nil
			}
		}
	}

	return resActions
}

// IsIAMSysInstance judge whether the resource type is a system instance in iam resource
func IsIAMSysInstance(resourceType TypeID) bool {
	return strings.HasPrefix(string(resourceType), IAMSysInstTypePrefix)
}

// IsCMDBSysInstance judge whether the resource type is a system instance in cmdb resource
func IsCMDBSysInstance(resourceType meta.ResourceType) bool {
	return strings.HasPrefix(string(resourceType), meta.CMDBSysInstTypePrefix)
}

// isIAMSysInstanceSelection judge whether the instance selection is a system instance selection in iam resource
func isIAMSysInstanceSelection(instanceSelectionID InstanceSelectionID) bool {
	return strings.Contains(string(instanceSelectionID), IAMSysInstTypePrefix)
}

// isIAMSysInstanceAction judge whether the action is a system instance action in iam resource
func isIAMSysInstanceAction(actionID ActionID) bool {
	return strings.Contains(string(actionID), IAMSysInstTypePrefix)
}

// GetModelIDFromIamSysInstance get model id from iam system instance
func GetModelIDFromIamSysInstance(resourceType TypeID) (int64, error) {
	if !IsIAMSysInstance(resourceType) {
		return 0, fmt.Errorf("resourceType %s is not an iam system instance, it must start with prefix %s",
			resourceType, IAMSysInstTypePrefix)
	}
	modelIDStr := strings.TrimPrefix(string(resourceType), IAMSysInstTypePrefix)
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		blog.ErrorJSON("modelID convert to int64 failed, err:%s, input:%s", err, modelID)
		return 0, fmt.Errorf("get model id failed, parse to int err:%s, the format of resourceType:%s is wrong",
			err.Error(), resourceType)
	}

	return modelID, nil
}

// GetActionTypeFromIAMSysInstance get action type from iam system instance
func GetActionTypeFromIAMSysInstance(actionID ActionID) ActionType {
	actionIDStr := string(actionID)
	return ActionType(actionIDStr[:strings.Index(actionIDStr, "_")])
}
