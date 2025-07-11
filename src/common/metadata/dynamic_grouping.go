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
	"errors"
	"fmt"
	"regexp"
	"time"

	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"

	"github.com/google/uuid"
)

// operators support for dynamic group conditions in system.
const (
	// DynamicGroupOperatorEQ eq operator.
	DynamicGroupOperatorEQ = "$eq"

	// DynamicGroupOperatorNE ne operator.
	DynamicGroupOperatorNE = "$ne"

	// DynamicGroupOperatorIN in operator.
	DynamicGroupOperatorIN = "$in"

	// DynamicGroupOperatorNIN nin operator.
	DynamicGroupOperatorNIN = "$nin"

	// DynamicGroupOperatorLTE lte operator.
	DynamicGroupOperatorLTE = "$lte"

	// DynamicGroupOperatorGTE gte operator.
	DynamicGroupOperatorGTE = "$gte"

	// DynamicGroupOperatorLIKE like operator.
	DynamicGroupOperatorLIKE = "$regex"
)

var (
	// DynamicGroupOperators all operators -> current newest operators.
	DynamicGroupOperators = map[string]string{
		DynamicGroupOperatorEQ:           DynamicGroupOperatorEQ,
		DynamicGroupOperatorNE:           DynamicGroupOperatorNE,
		DynamicGroupOperatorIN:           DynamicGroupOperatorIN,
		DynamicGroupOperatorNIN:          DynamicGroupOperatorNIN,
		DynamicGroupOperatorLTE:          DynamicGroupOperatorLTE,
		DynamicGroupOperatorGTE:          DynamicGroupOperatorGTE,
		DynamicGroupOperatorLIKE:         DynamicGroupOperatorLIKE,
		string(filter.Contains):          string(filter.Contains),
		string(filter.ContainsSensitive): string(filter.ContainsSensitive),
	}

	// DynamicGroupConditionTypes all condition object types of dynamic group.
	DynamicGroupConditionTypes = map[string]map[string]string{
		// host dynamic group.
		common.BKInnerObjIDHost: {
			common.BKInnerObjIDSet:    common.BKInnerObjIDSet,
			common.BKInnerObjIDModule: common.BKInnerObjIDModule,
			common.BKInnerObjIDHost:   common.BKInnerObjIDHost,
		},

		// set dynamic group.
		common.BKInnerObjIDSet: {
			common.BKInnerObjIDSet: common.BKInnerObjIDSet,
		},
	}
)

// Validatefunc is func callback for validating.
type Validatefunc func(objectID string) ([]Attribute, error)

// DynamicGroupCondition is target resource search condition on fields level.
type DynamicGroupCondition struct {
	// Field is target field name for index resource.
	Field string `json:"field" bson:"field"`

	// Operator is index operator type, eg $ne/$eq/$in/$nin.
	Operator string `json:"operator" bson:"operator"`

	// Value is target field value for index resource(integer or string).
	Value interface{} `json:"value" bson:"value"`
}

// Validate validates dynamic group conditions format.
func (c *DynamicGroupCondition) Validate(attributeMap map[string]string) error {
	operator, isSupport := DynamicGroupOperators[c.Operator]
	if !isSupport {
		return fmt.Errorf("not support operator, %s", c.Operator)
	}

	if c.Field == common.BKDefaultField {
		return nil
	}

	attributeType, isSupport := attributeMap[c.Field]
	if !isSupport {
		return fmt.Errorf("not support condition field, %+v", c.Field)
	}

	attrType, err := getAttributeType(attributeType)
	if err != nil {
		return err
	}

	if c.Field == common.BKServiceTemplateIDField {
		if c.Operator != DynamicGroupOperatorIN {
			return fmt.Errorf("service template field only support $in operator, not support operator, %s", c.Operator)
		}
	}

	switch attrType {
	case boolType:
		if operator != DynamicGroupOperatorEQ {
			return fmt.Errorf("bool type only support $eq operator, not support operator, %s", c.Operator)
		}
		return validAttributeValueType(attrType, c.Value)
	case dateType:
		if operator != DynamicGroupOperatorGTE && operator != DynamicGroupOperatorLTE {
			return fmt.Errorf("date type only support $gte & $lte operator, not support operator, %s", c.Operator)
		}
		return validAttributeValueType(attrType, c.Value)
	}

	switch operator {
	case DynamicGroupOperatorEQ, DynamicGroupOperatorNE:
		return validAttributeValueType(attrType, c.Value)
	case DynamicGroupOperatorIN, DynamicGroupOperatorNIN:
		valueArr, ok := c.Value.([]interface{})
		if !ok {
			return fmt.Errorf("operator %s only support array value, not support value, %+v", c.Operator, c.Value)
		}

		for _, value := range valueArr {
			if err := validAttributeValueType(attrType, value); err != nil {
				return err
			}
		}
	case DynamicGroupOperatorLTE, DynamicGroupOperatorGTE:
		if attrType != numericType {
			return fmt.Errorf("operator %s only support numeric value, not support attribute type, %s", c.Operator,
				attributeType)
		}
		return validAttributeValueType(attrType, c.Value)
	case DynamicGroupOperatorLIKE:
		if attrType != stringType {
			return fmt.Errorf("operator %s only support string value, not support attribute type, %s", c.Operator,
				attributeType)
		}

		return validAttributeValueType(attrType, c.Value)
	}

	return nil
}

// VerifyRegexValidity 验证正则表达式的合法性
func (c *DynamicGroupCondition) VerifyRegexValidity() error {
	// 验证 value 是否为空
	if c.Value == nil {
		blog.Errorf("HTTP request body data is not set, err: value not set, regex: %v", c.Value)
		return errors.New("value not set")
	}
	// 模糊匹配时需要验证正则表达式的合法性
	if c.Operator != common.BKDBLIKE {
		return nil
	}
	strValue := util.GetStrByInterface(c.Value)
	if strValue == "" {
		blog.Errorf("HTTP request body data is not set, err: value not set, regex: %v", c.Value)
		return errors.New("value not set")
	}
	if _, err := regexp.Compile(strValue); err != nil {
		blog.Errorf("the regular expression's type assertion failed, err: %v, regex: %v", err, c.Value)
		return err
	}
	return nil
}

func validAttributeValueType(attrType string, value interface{}) error {
	switch attrType {
	case stringType:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("attribute only support string value, not support value, %+v", value)
		}
	case numericType:
		if !util.IsNumeric(value) {
			return fmt.Errorf("attribute only support numeric value, not support value, %+v", value)
		}
	case boolType:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("attribute only support bool value, not support value, %+v", value)
		}
	case dateType:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("date attribute only support string value, not support value, %+v", value)
		}
	}

	return nil
}

const (
	numericType = "numeric"
	boolType    = "bool"
	stringType  = "string"
	dateType    = "date"
)

func getAttributeType(attributeType string) (string, error) {
	switch attributeType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeEnum, common.FieldTypeEnumMulti,
		common.FieldTypeTimeZone, common.FieldTypeUser, common.FieldTypeList:
		return stringType, nil
	case common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeOrganization, common.FieldTypeEnumQuote:
		return numericType, nil
	case common.FieldTypeBool:
		return boolType, nil
	case common.FieldTypeDate:
		return dateType, nil
	default:
		return "", fmt.Errorf("not support attribute type, %s", attributeType)
	}
}

// DynamicGroupInfoCondition is condition for dynamic grouping, user could search
// target source base on the conditions.
type DynamicGroupInfoCondition struct {
	// ObjID is cmdb object id, could be host/set now.
	ObjID string `json:"bk_obj_id" bson:"bk_obj_id"`

	// Condition is search condition on fields level.
	// Example: bk_host_name $eq my-host just index host which name is "my-host".
	Condition []DynamicGroupCondition `json:"condition" bson:"condition"`

	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty" bson:"time_condition,omitempty"`
}

// Validate validates dynamic group info conditions format.
func (c *DynamicGroupInfoCondition) Validate(validatefunc Validatefunc) error {
	attributes, err := validatefunc(c.ObjID)
	if err != nil {
		return fmt.Errorf("validate dynamic group failed, %+v", err)
	}

	attributeMap := make(map[string]string)
	for _, attribute := range attributes {
		attributeMap[attribute.PropertyID] = attribute.PropertyType
	}

	switch c.ObjID {
	case common.BKInnerObjIDSet:
		attributeMap[common.BKSetIDField] = common.FieldTypeInt

	case common.BKInnerObjIDModule:
		attributeMap[common.BKModuleIDField] = common.FieldTypeInt

	case common.BKInnerObjIDHost:
		attributeMap[common.BKHostIDField] = common.FieldTypeInt
		attributeMap[common.BKCloudIDField] = common.FieldTypeInt
	}

	blog.Infof("validate info conditions, object[%s] attributes[%+v]", c.ObjID, attributeMap)

	for _, cond := range c.Condition {
		if err := cond.Validate(attributeMap); err != nil {
			return err
		}
	}
	return nil
}

// DynamicGroupInfo is info field in DynamicGroup struct.
type DynamicGroupInfo struct {
	// Condition is dynamic group index lock conditions set.
	Condition []DynamicGroupInfoCondition `json:"condition" bson:"condition"`

	// VariableCondition is dynamic group index variable conditions set.
	VariableCondition []DynamicGroupInfoCondition `json:"variable_condition" bson:"variable_condition"`
}

// Validate validates dynamic group info format, it's OK if conditions empty in this level.
func (c *DynamicGroupInfo) Validate(objectID string, validatefunc Validatefunc) error {
	err := ValidDynamicGroupCond(c.Condition, objectID, validatefunc, make(map[string]map[string]struct{}))
	if err != nil {
		return err
	}

	err = ValidDynamicGroupCond(c.VariableCondition, objectID, validatefunc, GetMapFromDynamicCond(c.Condition))
	if err != nil {
		return err
	}

	return nil
}

// GetMapFromDynamicCond get map from dynamic group condition
func GetMapFromDynamicCond(condArr []DynamicGroupInfoCondition) map[string]map[string]struct{} {
	result := make(map[string]map[string]struct{})

	for _, cond := range condArr {
		if _, exist := result[cond.ObjID]; !exist {
			result[cond.ObjID] = map[string]struct{}{}
		}

		for _, item := range cond.Condition {
			result[cond.ObjID][item.Field] = struct{}{}
		}

		if cond.TimeCondition == nil {
			continue
		}

		for _, item := range cond.TimeCondition.Rules {
			result[cond.ObjID][item.Field] = struct{}{}
		}
	}

	return result
}

// ValidDynamicGroupCond validate dynamic group info condition
func ValidDynamicGroupCond(condition []DynamicGroupInfoCondition, objectID string, validatefunc Validatefunc,
	checkDupMap map[string]map[string]struct{}) error {

	types, isSupport := DynamicGroupConditionTypes[objectID]
	if !isSupport {
		return fmt.Errorf("not support dynamic group type, %s", objectID)
	}

	for _, cond := range condition {
		for _, item := range cond.Condition {
			if err := item.VerifyRegexValidity(); err != nil {
				blog.Errorf("verify regex validity failed, err: %v, input: %v, objectID: %s", err, item, objectID)
				return err
			}

			if _, ok := checkDupMap[cond.ObjID]; !ok {
				continue
			}

			if _, exist := checkDupMap[cond.ObjID][item.Field]; exist {
				return errors.New("cannot set the same field")
			}
		}

		if _, isSupport = types[cond.ObjID]; !isSupport {
			return fmt.Errorf("not support condition type[%s] for %s dynamic group", cond.ObjID, objectID)
		}

		if err := cond.Validate(validatefunc); err != nil {
			return err
		}

		if cond.TimeCondition == nil {
			continue
		}

		for _, rule := range cond.TimeCondition.Rules {
			if _, ok := checkDupMap[cond.ObjID]; !ok {
				continue
			}

			if _, exist := checkDupMap[cond.ObjID][rule.Field]; exist {
				return errors.New("cannot set the same field")
			}
		}
	}

	return nil
}

// DynamicGroup is dynamic grouping of conditions for host/set data searching.
type DynamicGroup struct {
	// AppID is application id which dynamic group belongs to.
	AppID int64 `json:"bk_biz_id" bson:"bk_biz_id"`

	// ID is dynamic group instance unique id.
	ID string `json:"id" bson:"id"`

	// Name is dynamic group name.
	Name string `json:"name" bson:"name"`

	// ObjID is cmdb object id, could be host/set now.
	ObjID string `json:"bk_obj_id" bson:"bk_obj_id"`

	// Info is dynamic group core conditions information.
	Info DynamicGroupInfo `json:"info" bson:"info"`

	// CreateUser create user name.
	CreateUser string `json:"create_user" bson:"create_user"`

	// ModifyUser modify user name.
	ModifyUser string `json:"modify_user" bson:"modify_user"`

	// CreateTime create timestamp.
	CreateTime time.Time `json:"create_time" bson:"create_time"`

	// UpdateTime last update timestamp.
	UpdateTime time.Time `json:"last_time" bson:"last_time"`
}

// Validate validates dynamic group format.
func (g *DynamicGroup) Validate(validatefunc Validatefunc) error {
	if g.AppID <= 0 {
		return errors.New("empty bk_biz_id")
	}

	if len(g.Name) == 0 {
		return errors.New("empty name")
	}

	// check object id.
	if len(g.ObjID) == 0 {
		return errors.New("empty bk_obj_id")
	}

	// check conditions format.
	if len(g.Info.Condition) == 0 && len(g.Info.VariableCondition) == 0 {
		// it's not OK if conditions empty in this level.
		return errors.New("info.condition and info.variable_condition can not be empty at the same time")
	}
	return g.Info.Validate(g.ObjID, validatefunc)
}

// DynamicGroupBatch is batch result struct of dynamic group.
type DynamicGroupBatch struct {
	// Count batch count.
	Count uint64 `json:"count"`

	// Info batch data.
	Info []DynamicGroup `json:"info"`
}

// SearchDynamicGroupResult is result struct for dynamic group searching action.
type SearchDynamicGroupResult struct {
	BaseResp `json:",inline"`
	Data     DynamicGroupBatch `json:"data"`
}

// GetDynamicGroupResult is result struct for dynamic group detail query action.
type GetDynamicGroupResult struct {
	BaseResp `json:",inline"`
	Data     DynamicGroup `json:"data"`
}

// NewDynamicGroupID creates and returns a new dynamic group string unique ID.
func NewDynamicGroupID() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// ExecuteOption execute dynamic group option
type ExecuteOption struct {
	VariableCondition []DynamicGroupInfoCondition `json:"variable_condition"`
	Fields            []string                    `json:"fields"`
	Page              BasePage                    `json:"page"`
	DisableCounter    bool                        `json:"disable_counter"`
}
