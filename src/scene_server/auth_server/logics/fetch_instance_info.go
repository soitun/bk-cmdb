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

package logics

import (
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	sdktypes "configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/scene_server/auth_server/types"
)

var resourceParentMap = iam.GetResourceParentMap()

// FetchInstanceInfo fetch resource instances' specified attributes info using instance ids
func (lgc *Logics) FetchInstanceInfo(kit *rest.Kit, resourceType iam.TypeID, filter *types.FetchInstanceInfoFilter,
	extraCond map[string]interface{}) ([]map[string]interface{}, error) {

	idField := GetResourceIDField(resourceType)
	nameField := GetResourceNameField(resourceType)
	if idField == "" || nameField == "" {
		blog.Errorf("request type %s is invalid, rid: %s", resourceType, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}

	if len(filter.Attrs) == 0 {
		return make([]map[string]interface{}, 0), nil
	}

	// if attribute filter is set, add id attribute and convert display_name to the real name field
	var attrs []string
	needPath := false
	if len(filter.Attrs) > 0 {
		attrs = append(filter.Attrs, idField)
		for index, attr := range attrs {
			if attr == types.NameField {
				attrs[index] = nameField
				continue
			}
			if attr == sdktypes.IamPathKey {
				needPath = true
			}
		}
		if needPath {
			for _, parent := range resourceParentMap[resourceType] {
				attrs = append(attrs, GetResourceIDField(parent))
			}
		}
	}

	cond := make(map[string]interface{})
	if isResourceIDStringType(resourceType) {
		cond[idField] = map[string]interface{}{common.BKDBIN: filter.IDs}
	} else {
		ids := make([]int64, len(filter.IDs))
		for idx, idStr := range filter.IDs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				blog.Errorf("id %s parse int failed, error: %s, rid: %s, skip it", idStr, err.Error(), kit.Rid)
				return nil, err
			}
			ids[idx] = id
		}
		cond[idField] = map[string]interface{}{
			common.BKDBIN: ids,
		}
	}

	if len(extraCond) > 0 {
		cond = map[string]interface{}{common.BKDBAND: []map[string]interface{}{cond, extraCond}}
	}

	param := metadata.PullResourceParam{Condition: cond, Fields: attrs, Limit: common.BKNoLimit, Offset: 0}
	instances, err := lgc.searchAuthResource(kit, param, resourceType)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}

	// covert id and display_name field
	for _, instance := range instances.Info {
		instance[types.IDField] = util.GetStrByInterface(instance[idField])
		if instance[nameField] != nil {
			instance[types.NameField] = util.GetStrByInterface(instance[nameField])
		}
		if needPath {
			instance[sdktypes.IamPathKey], err = lgc.getResourceIamPath(kit, resourceType, instance)
			if err != nil {
				blog.ErrorJSON("getResourceIamPath failed, error: %s, instance: %s, rid: %s", err.Error(), instance,
					kit.Rid)
				return nil, err
			}
		}
	}
	return instances.Info, nil
}

// FetchHostInfo fetch hosts' specified attributes info using host ids
func (lgc *Logics) FetchHostInfo(kit *rest.Kit, resourceType iam.TypeID, filter *types.FetchInstanceInfoFilter) (
	[]map[string]interface{}, error) {

	if resourceType != iam.Host {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField)
	}

	if len(filter.Attrs) == 0 {
		return make([]map[string]interface{}, 0), nil
	}

	// if attribute filter is set, add id attribute and convert display_name to the real name field
	var attrs []string
	needPath := false
	hasName := false
	if len(filter.Attrs) > 0 {
		attrs = append(filter.Attrs, common.BKHostIDField)
		for index, attr := range attrs {
			if attr == types.NameField {
				attrs[index] = common.BKHostInnerIPField
				hasName = true
				continue
			}
			if attr == sdktypes.IamPathKey {
				needPath = true
			}
		}
		if hasName {
			attrs = append(attrs, common.BKHostInnerIPv6Field)
			attrs = append(attrs, common.BKCloudIDField)
		}
	}

	hostIDs := make([]int64, len(filter.IDs))
	for idx, idStr := range filter.IDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			blog.Errorf("id %s parse int failed, error: %s, rid: %s, skip it", idStr, err.Error(), kit.Rid)
			return nil, err
		}
		hostIDs[idx] = id
	}

	hostIDLen := len(hostIDs)
	hosts := make([]map[string]interface{}, 0)
	for offset := 0; offset < hostIDLen; offset += 500 {
		limit := offset + 500
		if limit > hostIDLen {
			limit = hostIDLen
		}
		hostParam := &metadata.ListWithIDOption{
			IDs:    hostIDs[offset:limit],
			Fields: attrs,
		}
		hostArrStr, err := lgc.CoreAPI.CacheService().Cache().Host().ListHostWithHostID(kit.Ctx, kit.Header, hostParam)
		if err != nil {
			blog.Errorf("get hosts from cache failed, err: %v, hostIDs: %+v", err, hostIDs)
			return nil, err
		}

		hostArr := make([]map[string]interface{}, 0)
		err = json.Unmarshal([]byte(hostArrStr), &hostArr)
		if err != nil {
			blog.Errorf("unmarshal hosts %s failed, err: %v", hostArrStr, err)
			return nil, err
		}

		hosts = append(hosts, hostArr...)
	}

	if len(hosts) == 0 {
		return hosts, nil
	}

	return lgc.enrichHostInfo(kit, hosts, hasName, needPath)
}

func (lgc *Logics) enrichHostInfo(kit *rest.Kit, hosts []map[string]interface{}, hasName bool, needPath bool) (
	[]map[string]interface{}, error) {

	cnt := len(hosts)
	cloudIDList := make([]int64, cnt)
	hostIDList := make([]int64, cnt)

	for index, host := range hosts {
		if hasName {
			cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
			if err != nil {
				blog.Errorf("parse cloud area id failed, err: %v, host: %+v", err, host)
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
			}
			cloudIDList[index] = cloudID
		}

		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.Errorf("parse host id failed, err: %v, host: %+v", err, host)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}
		hostIDList[index] = hostID
	}

	// get cloud area for display name use
	var cloudMap map[int64]string
	var err error
	if hasName {
		cloudMap, err = lgc.getCloudNameMapByIDs(kit, cloudIDList)
		if err != nil {
			return nil, err
		}
	}

	var hostPathMap map[int64][]string
	if needPath {
		hostPathMap, err = lgc.getHostIamPath(kit, iam.Host, hostIDList)
		if err != nil {
			return nil, err
		}
	}

	// covert id and display_name field
	for _, host := range hosts {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.Errorf("parse host id failed, err: %v, host: %+v", err, host)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}
		// add id field
		host[types.IDField] = strconv.FormatInt(hostID, 10)

		if hasName {
			cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
			if err != nil {
				blog.Errorf("parse cloud area id failed, err: %v, host: %+v", err, host)
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
			}
			ip := util.GetStrByInterface(host[common.BKHostInnerIPField])
			ipv6 := util.GetStrByInterface(host[common.BKHostInnerIPv6Field])
			host[types.NameField] = metadata.GetHostDisplayName(ip, ipv6, cloudMap[cloudID])
		}

		if needPath {
			host[sdktypes.IamPathKey] = hostPathMap[hostID]
		}
	}

	return hosts, nil
}

// FetchObjInstInfo fetch object instances' specified attributes info using instance ids
func (lgc *Logics) FetchObjInstInfo(kit *rest.Kit, resourceType iam.TypeID, filter *types.FetchInstanceInfoFilter) (
	[]map[string]interface{}, error) {

	if !iam.IsIAMSysInstance(resourceType) {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField)
	}

	if len(filter.Attrs) == 0 {
		return make([]map[string]interface{}, 0), nil
	}

	// if attribute filter is set, add id attribute and convert display_name to the real name field
	var attrs []string
	needPath := false
	if len(filter.Attrs) > 0 {
		attrs = append(filter.Attrs, common.BKInstIDField)
		for index, attr := range attrs {
			if attr == types.NameField {
				attrs[index] = common.BKInstNameField
				continue
			}
			if attr == sdktypes.IamPathKey {
				needPath = true
			}
		}
		if needPath {
			for _, parent := range resourceParentMap[resourceType] {
				attrs = append(attrs, GetResourceIDField(parent))
			}
		}
	}

	ids := make([]int64, len(filter.IDs))
	for idx, idStr := range filter.IDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			blog.Errorf("id %s parse int failed, error: %s, rid: %s, skip it", idStr, err.Error(), kit.Rid)
			return nil, err
		}
		ids[idx] = id
	}

	instances := make([]map[string]interface{}, 0)
	query := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKInstIDField: map[string]interface{}{common.BKDBIN: ids},
		},
		Fields: attrs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	objID, err := lgc.GetObjIDFromResourceType(kit.Ctx, kit.Header, resourceType)

	if err != nil {
		blog.ErrorJSON("get object id from resource type failed, error: %s, resource type: %s, rid: %s",
			err, resourceType, kit.Rid)
		return nil, err
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, query)
	if err != nil {
		blog.Errorf("read object %s instances by ids(%+v) failed, err: %v, rid: %s", objID, ids, err, kit.Rid)
		return nil, err
	}

	// covert id and display_name field
	for _, instance := range result.Info {
		instance[types.IDField] = util.GetStrByInterface(instance[common.BKInstIDField])
		if instance[common.BKInstNameField] != nil {
			instance[types.NameField] = util.GetStrByInterface(instance[common.BKInstNameField])
		}
		if needPath {
			var err error
			instance[sdktypes.IamPathKey], err = lgc.getResourceIamPath(kit, resourceType, instance)
			if err != nil {
				blog.ErrorJSON("get iam path failed, err: %s, instance: %s, rid: %s", err, instance, kit.Rid)
				return nil, err
			}
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// ValidateFetchInstanceInfoRequest TODO
func (lgc *Logics) ValidateFetchInstanceInfoRequest(kit *rest.Kit,
	req *types.PullResourceReq) (*types.FetchInstanceInfoFilter, error) {
	filter, ok := req.Filter.(types.FetchInstanceInfoFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for fetch_instance_info method, rid: %s", req.Filter,
			kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}

	if len(filter.IDs) == 0 {
		blog.ErrorJSON("request filter %s ids not set for fetch_instance_info method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "filter.ids")
	}
	return &filter, nil
}

// getResourceIamPath TODO
// get resource iam path
func (lgc *Logics) getResourceIamPath(kit *rest.Kit, resourceType iam.TypeID,
	instance map[string]interface{}) ([]string, error) {
	if resourceType == iam.Host {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField)
	}

	iamPath := make([]string, 0)
	// currently all resources only have one layer TODO support multiple layers if needed
	for _, parent := range resourceParentMap[resourceType] {
		iamPath = append(iamPath,
			"/"+string(parent)+","+util.GetStrByInterface(instance[GetResourceIDField(parent)])+"/")
	}
	return iamPath, nil
}

func (lgc *Logics) getHostIamPath(kit *rest.Kit, resourceType iam.TypeID, hostList []int64) (map[int64][]string,
	error) {

	if resourceType != iam.Host {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField)
	}

	// get host iam path, either in resource pool directory or in business
	// TODO: support host in business module when topology is supported
	defaultBizID, err := lgc.GetResourcePoolBizID(kit)
	if err != nil {
		return nil, err
	}

	req := &metadata.HostModuleRelationRequest{
		HostIDArr: hostList,
		Fields:    []string{common.BKHostIDField, common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	res, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, req)
	if err != nil {
		blog.Errorf("GetHostModuleRelation by host id %v failed, err: %s, rid: %s", hostList, err, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		return make(map[int64][]string), nil
	}

	relationMap := make(map[int64][]string)
	for _, relation := range res.Info {
		var path string
		if relation.AppID == defaultBizID {
			path = "/" + string(iam.SysResourcePoolDirectory) + "," + strconv.FormatInt(relation.ModuleID, 10) + "/"
		} else {
			path = "/" + string(iam.Business) + "," + strconv.FormatInt(relation.AppID, 10) + "/"
		}

		if _, exist := relationMap[relation.HostID]; !exist {
			relationMap[relation.HostID] = make([]string, 0)
		}

		relationMap[relation.HostID] = append(relationMap[relation.HostID], path)

	}

	return relationMap, nil
}

// FetchSetModuleNameInfo fetch set & module resource name info for no permission apply url use
func (lgc *Logics) FetchSetModuleNameInfo(kit *rest.Kit, resType iam.TypeID, filter *types.FetchInstanceInfoFilter) (
	[]map[string]interface{}, error) {

	var objID string
	switch resType {
	case iam.Set:
		objID = common.BKInnerObjIDSet
	case iam.Module:
		objID = common.BKInnerObjIDModule
	default:
		blog.Errorf("resource type %s is invalid, rid: %s", resType, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}

	if len(filter.Attrs) == 0 && !util.Contains(filter.Attrs, types.NameField) {
		return make([]map[string]interface{}, 0), nil
	}

	filterIDs, err := util.SliceStrToInt64(filter.IDs)
	if err != nil {
		blog.Errorf("parse ids(%+v) to []int failed, err: %v, rid: %s", filter.IDs, err, kit.Rid)
		return nil, err
	}
	idField := common.GetInstIDField(objID)
	cond := mapstr.MapStr{
		idField: mapstr.MapStr{common.BKDBIN: filterIDs},
	}

	nameField := common.GetInstNameField(objID)
	param := metadata.PullResourceParam{Condition: cond, Fields: []string{idField, nameField}, Limit: common.BKNoLimit,
		Offset: 0}
	instances, err := lgc.searchAuthResource(kit, param, resType)
	if err != nil {
		blog.Errorf("search auth resource failed, err: %v, param: %+v, rid: %s", err, param, kit.Rid)
		return nil, err
	}

	for _, instance := range instances.Info {
		instance[types.IDField] = util.GetStrByInterface(instance[idField])
		if instance[nameField] != nil {
			instance[types.NameField] = util.GetStrByInterface(instance[nameField])
		}
	}
	return instances.Info, nil
}
