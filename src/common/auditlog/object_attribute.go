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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type objectAttributeAuditLog struct {
	audit
}

// BatchGenerateAuditLog generate log content about model attributes in batches, where
// the parameter attrs must be obtained by the caller
func (h *objectAttributeAuditLog) BatchGenerateAuditLog(parameter *generateAuditCommonParameter, objID string,
	attrs []metadata.Attribute) ([]metadata.AuditLog, error) {

	if len(attrs) == 0 {
		return []metadata.AuditLog{}, nil
	}

	kit := parameter.kit

	objName, err := h.getObjNameByObjID(kit, objID)
	if err != nil {
		return nil, err
	}
	logs := make([]metadata.AuditLog, len(attrs))
	for index, attr := range attrs {
		if attr.ObjectID != objID {
			return nil, kit.CCError.CCError(common.CCErrCommParamsIsInvalid)
		}

		logs[index] = metadata.AuditLog{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelAttributeRes,
			Action:       parameter.action,
			BusinessID:   attr.BizID,
			ResourceID:   attr.ID,
			ResourceName: attr.PropertyName,
			OperateFrom:  parameter.operateFrom,
			OperationDetail: &metadata.ModelAttrOpDetail{
				BkObjID:   objID,
				BkObjName: objName,
				BasicOpDetail: metadata.BasicOpDetail{
					Details: parameter.NewBasicContent(attr.ToMapStr()),
				},
			},
		}
	}
	return logs, nil
}

// GenerateAuditLog generate audit of model attribute, if data is nil, will auto get current model attribute data by id.
func (h *objectAttributeAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, id int64,
	data *metadata.Attribute) (*metadata.AuditLog, error) {
	kit := parameter.kit

	if data == nil {
		// get current model attribute data by id.
		query := mapstr.MapStr{metadata.AttributeFieldID: id}
		rsp, err := h.clientSet.Model().ReadModelAttrByCondition(kit.Ctx, kit.Header,
			&metadata.QueryCondition{Condition: query})
		if err != nil {
			blog.Errorf("generate audit log of model attribute failed, failed to read model attribute, err: %v, "+
				"rid: %s", err.Error(), kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("generate audit log of model attribute failed, failed to read model attribute, err: %s, "+
				"rid: %s", kit.CCError.CCError(common.CCErrorModelNotFound), kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}

		data = &rsp.Info[0]
	}

	objName, err := h.getObjNameByObjID(kit, data.ObjectID)
	if err != nil {
		return nil, err
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelAttributeRes,
		Action:       parameter.action,
		BusinessID:   data.BizID,
		ResourceID:   id,
		ResourceName: data.PropertyName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.ModelAttrOpDetail{
			BkObjID:   data.ObjectID,
			BkObjName: objName,
			BasicOpDetail: metadata.BasicOpDetail{
				Details: parameter.NewBasicContent(data.ToMapStr()),
			},
		},
	}, nil
}

// NewObjectAttributeAuditLog TODO
func NewObjectAttributeAuditLog(clientSet coreservice.CoreServiceClientInterface) *objectAttributeAuditLog {
	return &objectAttributeAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
