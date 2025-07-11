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

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type modelAttrUnique struct {
}

var forbiddenCreateUniqueObjList = []string{
	common.BKInnerObjIDProject,
}

// CreateModelAttrUnique create model attribute unique
func (m *modelAttrUnique) CreateModelAttrUnique(kit *rest.Kit, objID string, data metadata.CreateModelAttrUnique) (
	*metadata.CreateOneDataResult, error) {

	if util.InStrArr(forbiddenCreateUniqueObjList, objID) {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	id, err := m.createModelAttrUnique(kit, objID, data)
	if err != nil {
		return nil, err
	}
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, nil
}

// UpdateModelAttrUnique TODO
func (m *modelAttrUnique) UpdateModelAttrUnique(kit *rest.Kit, objID string, id uint64, data metadata.UpdateModelAttrUnique) (*metadata.UpdatedCount, error) {
	err := m.updateModelAttrUnique(kit, objID, id, data)
	if err != nil {
		return nil, err
	}
	return &metadata.UpdatedCount{Count: 1}, nil
}

// DeleteModelAttrUnique TODO
func (m *modelAttrUnique) DeleteModelAttrUnique(kit *rest.Kit, objID string, id uint64) (*metadata.DeletedCount, error) {
	err := m.deleteModelAttrUnique(kit, objID, id)
	if err != nil {
		return nil, err
	}
	return &metadata.DeletedCount{Count: 1}, nil
}

// SearchModelAttrUnique TODO
func (m *modelAttrUnique) SearchModelAttrUnique(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryUniqueResult, error) {

	uniqueItems, err := m.searchModelAttrUnique(kit, inputParam)
	if nil != err {
		return &metadata.QueryUniqueResult{Info: []metadata.ObjectUnique{}}, err
	}
	dataResult := &metadata.QueryUniqueResult{Info: []metadata.ObjectUnique{}}
	dataResult.Count, err = m.countModelAttrUnique(kit, inputParam.Condition)
	if nil != err {
		return &metadata.QueryUniqueResult{Info: []metadata.ObjectUnique{}}, err
	}
	if len(uniqueItems) > 0 {
		dataResult.Info = uniqueItems
	}

	return dataResult, nil
}
