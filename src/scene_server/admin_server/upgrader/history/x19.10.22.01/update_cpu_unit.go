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

package x19_10_22_01

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// UpdateCpuUnit TODO
func UpdateCpuUnit(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: "bk_cpu_mhz",
	}
	doc := map[string]interface{}{
		"unit": "MHz",
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("UpdateCpuUnit failed, err: %+v", err)
		return fmt.Errorf("UpdateCpuUnit failed, err: %v", err)
	}
	return nil
}
