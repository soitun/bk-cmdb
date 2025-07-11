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

package metadata

import (
	"configcenter/src/common/mapstr"
)

const (
	// ClassificationFieldID TODO
	ClassificationFieldID = "id"
	// ClassFieldClassificationID TODO
	ClassFieldClassificationID = "bk_classification_id"
	// ClassFieldClassificationName TODO
	ClassFieldClassificationName = "bk_classification_name"
)
const (
	// ClassificationHostManageID TODO
	ClassificationHostManageID = "bk_host_manage"
	// ClassificationBizTopoID TODO
	ClassificationBizTopoID = "bk_biz_topo"
	// ClassificationOrganizationID TODO
	ClassificationOrganizationID = "bk_organization"
	// ClassificationNetworkID TODO
	ClassificationNetworkID = "bk_network"
	// ClassificationUncategorizedID TODO
	ClassificationUncategorizedID = "bk_uncategorized"
	// ClassificationTableID built-in table classification
	ClassificationTableID = "bk_table_classification"
)

// Model group classification initialization value
const (
	ClassificationHostManage    = "主机管理"
	ClassificationTopo          = "业务拓扑"
	ClassificationOrganization  = "组织架构"
	ClassificationNet           = "网络"
	ClassificationUncategorized = "未分类"
)

// Classification the classification metadata definition
type Classification struct {
	ID                 int64  `field:"id" json:"id" bson:"id" mapstructure:"id"     `
	ClassificationID   string `field:"bk_classification_id"  json:"bk_classification_id" bson:"bk_classification_id" mapstructure:"bk_classification_id"             `
	ClassificationName string `field:"bk_classification_name" json:"bk_classification_name" bson:"bk_classification_name" mapstructure:"bk_classification_name"`
	ClassificationType string `field:"bk_classification_type" json:"bk_classification_type" bson:"bk_classification_type" mapstructure:"bk_classification_type"`
	ClassificationIcon string `field:"bk_classification_icon" json:"bk_classification_icon" bson:"bk_classification_icon" mapstructure:"bk_classification_icon"`
	OwnerID            string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account" `
}

// Parse load the data from mapstr classification into classification instance
func (cli *Classification) Parse(data mapstr.MapStr) (*Classification, error) {

	err := mapstr.SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}
