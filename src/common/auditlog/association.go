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

package auditlog

import (
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

// instanceAssociationAuditLog provides methods to generate and save instance association audit log.
type instanceAssociationAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of instance association,
// if data is nil, will auto get data by id and instance association.
func (a *instanceAssociationAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, id int64, objID string,
	data *metadata.InstAsst) (*metadata.AuditLog, error) {
	kit := parameter.kit

	if data == nil {
		cond := metadata.InstAsstQueryCondition{
			Cond: metadata.QueryCondition{
				Condition: map[string]interface{}{metadata.AssociationFieldAssociationId: id},
			},
			ObjID: objID,
		}

		result, err := a.clientSet.Association().ReadInstAssociation(kit.Ctx, kit.Header, &cond)
		if err != nil {
			blog.Errorf("generate inst asst audit log failed, get instance association failed, err: %v, rid: %s",
				err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
		}

		if len(result.Info) == 0 || len(result.Info) > 1 {
			blog.Errorf("generate inst asst audit log failed, get instance association by id(%d) get none or "+
				"multiple result, rid: %s", id, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
		}

		data = &result.Info[0]
	}

	srcInstName, err := a.getInstNameByID(kit, data.ObjectID, data.InstID)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
	}

	targetInstName, err := a.getInstNameByID(kit, data.AsstObjectID, data.AsstInstID)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelInstanceType,
		ResourceType: metadata.InstanceAssociationRes,
		Action:       parameter.action,
		ResourceID:   data.InstID,
		ResourceName: srcInstName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.InstanceAssociationOpDetail{
			AssociationOpDetail: metadata.AssociationOpDetail{
				AssociationID:   data.ObjectAsstID,
				AssociationKind: data.AssociationKindID,
			},
			SourceModelID:      data.ObjectID,
			TargetModelID:      data.AsstObjectID,
			TargetInstanceID:   data.AsstInstID,
			TargetInstanceName: targetInstName,
		},
	}, nil
}

// NewInstanceAssociationAudit TODO
func NewInstanceAssociationAudit(clientSet coreservice.CoreServiceClientInterface) *instanceAssociationAuditLog {
	return &instanceAssociationAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
