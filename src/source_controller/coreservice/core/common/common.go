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

// Package common TODO
package common

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

var _ core.CommonOperation = (*commonOperation)(nil)

type commonOperation struct {
}

// New create a new instance manager instance
func New() core.CommonOperation {
	return &commonOperation{}
}

// GetDistinctField TODO
func (c *commonOperation) GetDistinctField(kit *rest.Kit, option *metadata.DistinctFieldOption) ([]interface{},
	errors.CCErrorCoder) {

	ret, err := mongodb.Client().Table(option.TableName).Distinct(kit.Ctx, option.Field, option.Filter)
	if err != nil {
		blog.Errorf("get distinct field failed, err: %v, option:%#v, rid: %s", err, *option, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return ret, nil
}

// GetDistinctCount 根据条件获取指定表中满足条件数据的数量
func (c *commonOperation) GetDistinctCount(kit *rest.Kit, option *metadata.DistinctFieldOption) (int64,
	errors.CCErrorCoder) {
	var count int64
	ret, err := mongodb.Client().Table(option.TableName).Distinct(kit.Ctx, option.Field, option.Filter)
	if err != nil {
		blog.Errorf("get distinct count failed, err: %v, option:%#v, rid: %s", err, *option, kit.Rid)
		return count, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	count = int64(len(ret))
	return count, nil
}
