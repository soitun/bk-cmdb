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
	"fmt"

	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
)

type syncTaskAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of cloud sync task, if data is nil, will auto get data by taskID.
func (h *syncTaskAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, taskID int64,
	data *metadata.CloudSyncTask) (
	*metadata.AuditLog, error) {
	kit := parameter.kit

	if data == nil {
		// get data by taskID.
		option := metadata.SearchCloudOption{
			Condition: mapstr.MapStr{common.BKCloudSyncTaskID: taskID},
		}

		res, err := h.clientSet.Cloud().SearchSyncTask(kit.Ctx, kit.Header, &option)
		if err != nil {
			blog.Errorf("generate audit log of cloud sync task, failed to read cloud sync task, err: %v, rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}
		if len(res.Info) <= 0 {
			blog.Errorf("generate audit log of cloud sync task failed, not find cloud sync task, rid: %s",
				kit.Rid)
			return nil, fmt.Errorf("generate audit log of cloud sync task failed, not find cloud sync task")
		}

		data = &res.Info[0]
	}

	dataMap, err := mapstruct.Struct2Map(data)
	if err != nil {
		blog.Errorf("convert cloud sync task(%+v) to map failed, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	return &metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudSyncTaskRes,
		Action:       parameter.action,
		ResourceID:   taskID,
		ResourceName: data.TaskName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.BasicOpDetail{
			Details: parameter.NewBasicContent(dataMap),
		},
	}, nil
}

// NewSyncTaskAuditLog TODO
func NewSyncTaskAuditLog(clientSet coreservice.CoreServiceClientInterface) *syncTaskAuditLog {
	return &syncTaskAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
