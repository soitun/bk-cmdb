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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/storage/driver/mongodb"
)

func (g *modelAttributeGroup) count(kit *rest.Kit, cond universalsql.Condition) (count int64, err error) {

	iCount, err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Find(cond.ToMapStr()).Count(kit.Ctx)
	return int64(iCount), err
}

func (g *modelAttributeGroup) save(kit *rest.Kit, group metadata.Group) (uint64, error) {

	id, err := mongodb.Client().NextSequence(kit.Ctx, common.BKTableNamePropertyGroup)
	if err != nil {
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	group.ID = int64(id)
	group.OwnerID = kit.SupplierAccount

	err = mongodb.Client().Table(common.BKTableNamePropertyGroup).Insert(kit.Ctx, group)
	return id, err
}

func (g *modelAttributeGroup) delete(kit *rest.Kit, cond universalsql.Condition) (uint64, error) {
	cnt, err := mongodb.Client().Table(common.BKTableNamePropertyGroup).DeleteMany(kit.Ctx, cond.ToMapStr())
	if err != nil {
		blog.ErrorJSON("delete model attribute group error. cond: %s, err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return 0, err
	}
	return cnt, nil
}

func (g *modelAttributeGroup) search(kit *rest.Kit, cond universalsql.Condition) ([]metadata.Group, error) {

	dataResult := make([]metadata.Group, 0)
	err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Find(cond.ToMapStr()).All(kit.Ctx, &dataResult)
	return dataResult, err
}

func (g *modelAttributeGroup) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (uint64, error) {

	cnt, err := mongodb.Client().Table(common.BKTableNamePropertyGroup).UpdateMany(kit.Ctx, cond.ToMapStr(), data)
	if err != nil {
		blog.ErrorJSON("update model attribute group error. cond: %s, err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return 0, err
	}
	return cnt, nil
}
