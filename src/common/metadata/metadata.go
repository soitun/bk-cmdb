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

// Package metadata TODO
package metadata

import (
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// ModelKind TODO
type ModelKind string

const (
	// LabelTrue and LabelFalse is used to define if the label value is a bool value,
	// which is true or false.
	LabelTrue string = "true"
	// LabelFalse TODO
	LabelFalse string = "false"

	// LabelBusinessID is used to define a label key which value is the business number value.
	LabelBusinessID string = "bk_biz_id"

	// LabelModelKind is used to define a label key which describe what kind of this object is.
	// this label key is always used with the value of PublicModelKindValue or BusinessModelKindValue.
	LabelModelKind string = "object_kind"
	// PublicModelKindValue TODO
	PublicModelKindValue ModelKind = "public_object"
	// BusinessModelKindValue TODO
	BusinessModelKindValue ModelKind = "business_object"
)

var (
	// LabelKeyNotExistError TODO
	LabelKeyNotExistError = errors.New("label key does not exist")
)

// MetadataBizField TODO
var MetadataBizField = "metadata.label.bk_biz_id"

// BizLabelNotExist TODO
var BizLabelNotExist = mapstr.MapStr{"metadata.label.bk_biz_id": mapstr.MapStr{"$exists": false}}

/*
func PublicAndBizCondition(meta Metadata) mapstr.MapStr {
	exist, bizID := meta.Label.Get(LabelBusinessID)
	if false == exist {
		bizID = ""
	}
	condArr := make([]mapstr.MapStr, 0)
	condArr = append(condArr, BizLabelNotExist, mapstr.MapStr{"metadata.label.bk_biz_id": bizID})
	return mapstr.MapStr{"$or": condArr}
}
*/

// BizIDFromMetadata TODO
func BizIDFromMetadata(meta Metadata) (int64, error) {
	var businessID int64
	var err error
	exist, bizID := meta.Label.Get(LabelBusinessID)
	if false == exist {
		return 0, nil
	}
	businessID, err = strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("convert business id failed, err: %+v", err)
	}
	return businessID, nil
}

// PublicAndBizCondition TODO
func PublicAndBizCondition(meta Metadata) mapstr.MapStr {
	var businessID int64
	var err error
	exist, bizID := meta.Label.Get(LabelBusinessID)
	if false == exist {
		return NewPublicOrBizConditionByBizID(0)
	}

	bizIDStr := util.GetStrByInterface(bizID)
	if len(bizIDStr) > 0 {
		businessID, err = strconv.ParseInt(bizIDStr, 10, 64)
		if err != nil {
			blog.Errorf("PublicAndBizCondition parse business id failed, generate public condition only, bizID: %+v, err: %+v",
				bizID, err)
			businessID = 0
		}
	}
	return NewPublicOrBizConditionByBizID(businessID)
}

// NewPublicOrBizConditionByBizID new a query condition
func NewPublicOrBizConditionByBizID(businessID int64) mapstr.MapStr {
	condArr := make([]mapstr.MapStr, 0)
	condArr = append(condArr, BizLabelNotExist)
	if businessID != 0 {
		condArr = append(condArr, mapstr.MapStr{"metadata.label.bk_biz_id": strconv.FormatInt(businessID, 10)})
	}
	return mapstr.MapStr{common.BKDBOR: condArr}
}

const (
	// BKMetadata TODO
	BKMetadata string = "metadata"
	// BKLabel TODO
	BKLabel string = "label"
)

// Label define the Label type used to limit the scope of application of resources
type Label map[string]string

// NewMetaDataFromBusinessID TODO
func NewMetaDataFromBusinessID(value string) Metadata {
	label := make(Label)
	label[LabelBusinessID] = value
	meta := Metadata{Label: label}
	return meta
}

// NewMetadata TODO
func NewMetadata(bizID int64) Metadata {
	return NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
}

// GetBusinessIDFromMeta TODO
func GetBusinessIDFromMeta(data interface{}) string {
	if nil == data {
		return ""
	}
	tmp, ok := data.(map[string]interface{})
	if !ok {
		return ""
	}
	label, ok := tmp[BKLabel].(map[string]interface{})
	if !ok {
		return ""
	}
	bizID, ok := label[LabelBusinessID].(string)
	if !ok {
		return ""
	}
	return bizID
}

// ParseBizIDFromData parse biz id from metadata, 0 for global case
// nil ==> 0, error
// [] ==> 0, error
// {}  ==> 0, nil
// {"label": {}} ==> 0, nil
// {"label": []} ==> 0, error
// {"label": {"bk_biz_id": 1}} ==> 1, nil
// {"label": {"bk_biz_id": "a"}} ==> 0, error
func ParseBizIDFromData(rawData mapstr.MapStr) (int64, error) {
	rawMetadata, exist := rawData.Get(BKMetadata)
	if exist == false {
		return 0, fmt.Errorf("invalid input, metadata field not exist")
	}
	js, _ := json.Marshal(rawMetadata)
	meta := new(Metadata)
	if err := json.Unmarshal(js, meta); err != nil {
		return 0, err
	}

	rawBizID, existed := meta.Label[LabelBusinessID]
	if !existed {
		return 0, nil
	}
	bizID, err := strconv.ParseInt(rawBizID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid biz id value, parse int failed, id: %+v, err: %+v", rawBizID, err)
	}
	return bizID, nil

}

// MetadataWrapper TODO
type MetadataWrapper struct {
	Metadata Metadata `json:"metadata"`
}

// Metadata  used to define the metadata for the resources
type Metadata struct {
	Label Label `field:"label" json:"label" bson:"label"`
}

// ParseBizID TODO
func (md *Metadata) ParseBizID() (int64, error) {
	if md == nil {
		return 0, errors.New("invalid nil matadata")
	}
	bizID, err := BizIDFromMetadata(*md)
	if err != nil {
		return 0, err
	}
	return bizID, nil
}

// ToMapStr TODO
func (md *Metadata) ToMapStr() mapstr.MapStr {
	return mapstr.MapStr{"label": md.Label.ToMapStr()}
}

// Set TODO
func (label Label) Set(key, value string) {
	label[key] = value
}

// Get TODO
func (label Label) Get(key string) (exist bool, value string) {
	value, exist = label[key]
	return
}

// ToMapStr TODO
func (label Label) ToMapStr() mapstr.MapStr {
	result := make(map[string]interface{})
	for key, value := range label {
		result[key] = value
	}
	return result
}

// IsTrue is used to check if the label key is a true value or not.
// if the key does not exist, it will return with a LabelKeyNotExistError,
// which can be used to check it, if you do care about it.
func (label Label) IsTrue(key string) (bool, error) {
	value, exist := label[key]
	if !exist {
		return false, LabelKeyNotExistError
	}

	return strconv.ParseBool(value)
}

// Int64 is used to get a int64 value from a label key.
// if the key does not exist, it will return with a LabelKeyNotExistError,
// which can be used to check it, if you do care about it.
func (label Label) Int64(key string) (int64, error) {
	value, exist := label[key]
	if !exist {
		return 0, LabelKeyNotExistError
	}

	return strconv.ParseInt(value, 10, 64)
}

// GetBusinessID TODO
func (label Label) GetBusinessID() (int64, error) {
	return label.Int64(LabelBusinessID)
}

// SetBusinessID TODO
func (label Label) SetBusinessID(id int64) {
	label[LabelBusinessID] = strconv.FormatInt(id, 10)
}

// GetModelKind TODO
func (label Label) GetModelKind() (ModelKind, error) {
	kind, exist := label[LabelModelKind]
	if !exist {
		return "", LabelKeyNotExistError
	}

	switch ModelKind(kind) {
	case PublicModelKindValue, BusinessModelKindValue:
		return ModelKind(kind), nil
	default:
		return "", fmt.Errorf("unknown model kind %s", kind)
	}
}

// SetModelKind TODO
func (label Label) SetModelKind(kind ModelKind) {
	label[LabelModelKind] = string(kind)
}
