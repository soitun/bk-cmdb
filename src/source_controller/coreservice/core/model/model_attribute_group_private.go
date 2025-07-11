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
	"configcenter/src/common/universalsql/mongo"
)

func (g *modelAttributeGroup) groupIDIsExists(kit *rest.Kit, objID, groupID string, modelBizID int64) (oneResult metadata.Group, isExists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldGroupID, Val: groupID})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldSupplierAccount, Val: kit.SupplierAccount})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})
	if modelBizID > 0 {
		cond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: modelBizID})
	}

	groups, err := g.search(kit, cond)
	if nil != err {
		return oneResult, isExists, err
	}

	if 0 != len(groups) {
		return groups[0], true, nil
	}

	return oneResult, isExists, nil
}

func (g *modelAttributeGroup) groupNameIsExists(kit *rest.Kit, objID, groupName string, modelBizID int64) (
	metadata.Group, bool, error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldGroupName, Val: groupName})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldSupplierAccount, Val: kit.SupplierAccount})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})

	if modelBizID > 0 {
		cond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: 0})
	}

	groups, err := g.search(kit, cond)
	if err != nil {
		return metadata.Group{}, false, err
	}

	if len(groups) != 0 {
		return groups[0], true, nil
	}

	return metadata.Group{}, false, nil
}

func (g *modelAttributeGroup) hasAttributes(kit *rest.Kit, objID string, groupIDS []string) (isExists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldSupplierAccount, Val: kit.SupplierAccount})
	cond.Element(&mongo.In{Key: metadata.AttributeFieldPropertyGroup, Val: groupIDS})

	attrs, err := g.model.SearchModelAttributes(kit, objID, metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	})

	if nil != err {
		return false, err
	}

	return 0 != attrs.Count, nil
}
