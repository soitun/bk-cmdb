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

// GenerateCommonActions generate all the common actions registered to IAM.
func GenerateCommonActions() []CommonAction {
	return []CommonAction{
		{
			Name:        "业务运维",
			EnglishName: "Business Maintainer",
			Actions: []ActionWithID{
				{ID: ViewBusinessResource},
				{ID: EditBusinessHost},
				{ID: BusinessHostTransferToResourcePool},
				{ID: CreateBusinessTopology},
				{ID: EditBusinessTopology},
				{ID: DeleteBusinessTopology},
				{ID: CreateBusinessServiceInstance},
				{ID: EditBusinessServiceInstance},
				{ID: DeleteBusinessServiceInstance},
				{ID: CreateBusinessServiceTemplate},
				{ID: EditBusinessServiceTemplate},
				{ID: DeleteBusinessServiceTemplate},
				{ID: CreateBusinessSetTemplate},
				{ID: EditBusinessSetTemplate},
				{ID: DeleteBusinessSetTemplate},
				{ID: CreateBusinessServiceCategory},
				{ID: EditBusinessServiceCategory},
				{ID: DeleteBusinessServiceCategory},
				{ID: CreateBusinessCustomQuery},
				{ID: EditBusinessCustomQuery},
				{ID: DeleteBusinessCustomQuery},
				{ID: EditBusinessCustomField},
				{ID: EditBusinessHostApply},
				{ID: FindBusiness},
			},
		},
		{
			Name:        "业务只读",
			EnglishName: "Business Visitor",
			Actions: []ActionWithID{
				{ID: ViewBusinessResource},
				{ID: FindBusiness},
			},
		},
		{
			Name:        "业务集运维",
			EnglishName: "Biz-set Maintainer",
			Actions: []ActionWithID{
				{ID: AccessBizSet},
				{ID: DeleteBizSet},
				{ID: ViewBizSet},
			},
		},
		{
			Name:        "业务集只读",
			EnglishName: "Biz-set Visitor",
			Actions: []ActionWithID{
				{ID: AccessBizSet},
				{ID: ViewBizSet},
			},
		},
		{
			Name:        "主机资源管理员",
			EnglishName: "Host Maintainer",
			Actions: []ActionWithID{
				{ID: ViewResourcePoolHost},
				{ID: CreateResourcePoolHost},
				{ID: EditResourcePoolHost},
				{ID: DeleteResourcePoolHost},
				{ID: ResourcePoolHostTransferToBusiness},
				{ID: ResourcePoolHostTransferToDirectory},
				{ID: ManageHostAgentID},
				{ID: CreateResourcePoolDirectory},
				{ID: EditResourcePoolDirectory},
				{ID: DeleteResourcePoolDirectory},
				{ID: CreateCloudAccount},
				{ID: EditCloudAccount},
				{ID: DeleteCloudAccount},
				{ID: FindCloudAccount},
				{ID: CreateCloudResourceTask},
				{ID: EditCloudResourceTask},
				{ID: DeleteCloudResourceTask},
				{ID: FindCloudResourceTask},
			},
		},
		{
			Name:        "开发者",
			EnglishName: "Developer",
			Actions: []ActionWithID{
				{ID: WatchHostEvent},
				{ID: WatchHostRelationEvent},
				{ID: WatchBizEvent},
				{ID: WatchSetEvent},
				{ID: WatchModuleEvent},
				{ID: WatchProcessEvent},
				{ID: WatchCommonInstanceEvent},
			},
		},
		{
			Name:        "模型关系维护人",
			EnglishName: "Model Maintainer",
			Actions: []ActionWithID{
				{ID: CreateModelGroup},
				{ID: EditModelGroup},
				{ID: DeleteModelGroup},
				{ID: EditBusinessLayer},
				{ID: EditModelTopologyView},
				{ID: CreateSysModel},
				{ID: EditSysModel},
				{ID: DeleteSysModel},
				{ID: CreateAssociationType},
				{ID: EditAssociationType},
				{ID: DeleteAssociationType},
			},
		},
		{
			Name:        "审计员",
			EnglishName: "Auditor",
			Actions: []ActionWithID{
				{ID: FindOperationStatistic},
				{ID: EditOperationStatistic},
				{ID: FindAuditLog},
			},
		},
	}
}
