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

package extensions

import (
	"strconv"

	"configcenter/src/ac"
	"configcenter/src/ac/iam"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// AuthManager TODO
type AuthManager struct {
	clientSet  apimachinery.ClientSetInterface
	Authorizer ac.AuthorizeInterface
	Viewer     ac.Viewer

	RegisterModuleEnabled        bool
	RegisterSetEnabled           bool
	RegisterAuditCategoryEnabled bool
	SkipReadAuthorization        bool
}

// NewAuthManager TODO
func NewAuthManager(clientSet apimachinery.ClientSetInterface, iamCli *iam.IAM) *AuthManager {
	return &AuthManager{
		clientSet:                    clientSet,
		Authorizer:                   iam.NewAuthorizer(clientSet),
		Viewer:                       iam.NewViewer(clientSet, iamCli),
		RegisterModuleEnabled:        false,
		RegisterSetEnabled:           false,
		SkipReadAuthorization:        true,
		RegisterAuditCategoryEnabled: false,
	}
}

// InstanceSimplify TODO
type InstanceSimplify struct {
	InstanceID int64  `field:"bk_inst_id" json:"bk_inst_id" bson:"bk_inst_id"`
	Name       string `field:"bk_inst_name" json:"bk_inst_name" bson:"bk_inst_name"`
	BizID      int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (is *InstanceSimplify) Parse(data mapstr.MapStr) (*InstanceSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if err != nil {
		return nil, err
	}

	bizID, err := ParseBizID(data)
	is.BizID = bizID
	return is, err
}

// ParseBizID TODO
func ParseBizID(data mapstr.MapStr) (int64, error) {
	/*
		data:
		{
		  "metadata": {
			"label": {
			  "bk_biz_id": "2"
			}
		  }
		}
		==> 2, nil
	*/

	metaInterface, exist := data[common.MetadataField]
	if !exist {
		return 0, nil
	}
	metaValue, ok := metaInterface.(map[string]interface{})
	if !ok {
		return 0, nil
	}
	return ParseBizIDFromMetadata(metaValue)
}

// ParseBizIDFromMetadata TODO
func ParseBizIDFromMetadata(metaValue map[string]interface{}) (int64, error) {
	labelInterface, exist := metaValue["label"]
	if !exist {
		return 0, nil
	}

	labelValue, ok := labelInterface.(map[string]interface{})
	if !ok {
		return 0, nil
	}

	bizID, exist := labelValue[common.BKAppIDField]
	if !exist {
		// 自定义层级的metadata.label中没有 bk_biz_id 字段
		return 0, nil
	}

	// if metadata biz id is of string type, convert it to int64, otherwise, convert the integer type of it to int64
	if bizIDStr, ok := bizID.(string); ok {
		return strconv.ParseInt(bizIDStr, 10, 64)
	}

	return util.GetInt64ByInterface(bizID)
}

// BusinessSimplify TODO
type BusinessSimplify struct {
	BKAppIDField   int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	BKAppNameField string `field:"bk_biz_name" json:"bk_biz_name" bson:"bk_biz_name"`
	BKOwnerIDField string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	IsDefault      int64  `field:"default" json:"default" bson:"default"`

	Maintainer string `field:"bk_biz_maintainer" json:"bk_biz_maintainer" bson:"bk_biz_maintainer"`
	Producer   string `field:"bk_biz_productor" json:"bk_biz_productor" bson:"bk_biz_productor"`
	Tester     string `field:"bk_biz_tester" json:"bk_biz_tester" bson:"bk_biz_tester"`
	Developer  string `field:"bk_biz_developer" json:"bk_biz_developer" bson:"bk_biz_developer"`
	Operator   string `field:"operator" json:"operator" bson:"operator"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (business *BusinessSimplify) Parse(data mapstr.MapStr) (*BusinessSimplify, error) {

	err := mapstr.SetValueToStructByTags(business, data)
	if err != nil {
		return nil, err
	}

	return business, err
}

// BizSetSimplify biz set simplify
type BizSetSimplify struct {
	BKBizSetIDField   int64  `field:"bk_biz_set_id"`
	BKBizSetNameField string `field:"bk_biz_set_name"`
}

// SetSimplify TODO
type SetSimplify struct {
	BKAppIDField   int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	BKSetIDField   int64  `field:"bk_set_id" json:"bk_set_id" bson:"bk_set_id"`
	BKSetNameField string `field:"bk_set_name" json:"bk_set_name" bson:"bk_set_name"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (is *SetSimplify) Parse(data mapstr.MapStr) (*SetSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if err != nil {
		return nil, err
	}

	return is, err
}

// ModuleSimplify TODO
type ModuleSimplify struct {
	BKAppIDField      int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	BKModuleIDField   int64  `field:"bk_module_id" json:"bk_module_id" bson:"bk_module_id"`
	BKModuleNameField string `field:"bk_module_name" json:"bk_module_name" bson:"bk_module_name"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (is *ModuleSimplify) Parse(data mapstr.MapStr) (*ModuleSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if err != nil {
		return nil, err
	}

	return is, err
}

// HostSimplify host simplify struct
type HostSimplify struct {
	BKAppIDField         int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	BKModuleIDField      int64  `field:"bk_module_id" json:"bk_module_id" bson:"bk_module_id"`
	BKSetIDField         int64  `field:"bk_set_id" json:"bk_set_id" bson:"bk_set_id"`
	BKHostIDField        int64  `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id"`
	BKHostNameField      string `field:"bk_host_name" json:"bk_host_name" bson:"bk_host_name"`
	BKHostInnerIPField   string `field:"bk_host_innerip" json:"bk_host_innerip" bson:"bk_host_innerip"`
	BKHostInnerIPv6Field string `field:"bk_host_innerip_v6" json:"bk_host_innerip_v6" bson:"bk_host_innerip_v6"`
	BKCloudID            int64  `field:"bk_cloud_id" json:"bk_cloud_id" bson:"bk_cloud_id"`
}

// Parse TODO
func (is *HostSimplify) Parse(data mapstr.MapStr) (*HostSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if err != nil {
		return nil, err
	}

	return is, err
}

// PlatSimplify TODO
type PlatSimplify struct {
	BKCloudIDField   int64  `field:"bk_cloud_id" json:"bk_cloud_id" bson:"bk_cloud_id"`
	BKCloudNameField string `field:"bk_cloud_name" json:"bk_cloud_name" bson:"bk_cloud_name"`
}

// Parse TODO
func (is *PlatSimplify) Parse(data mapstr.MapStr) (*PlatSimplify, error) {

	err := mapstr.SetValueToStructByTags(is, data)
	if err != nil {
		return nil, err
	}

	return is, err
}

// ProcessSimplify TODO
type ProcessSimplify struct {
	ProcessID    int64  `field:"bk_process_id" json:"bk_process_id" bson:"bk_process_id"`
	ProcessName  string `field:"bk_process_name" json:"bk_process_name" bson:"bk_process_name"`
	BKAppIDField int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
}

// Parse TODO
func (is *ProcessSimplify) Parse(data mapstr.MapStr) (*ProcessSimplify, error) {
	err := mapstr.SetValueToStructByTags(is, data)
	if err != nil {
		return nil, err
	}

	return is, err
}
