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

package x18_09_30_01

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func cleanBKCloud(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	clouds := []map[string]interface{}{}

	err = db.Table(common.BKTableNameBasePlat).Find(mapstr.New()).Sort("create_time").All(ctx, &clouds) // db.GetMutilByCondition(common.BKTableNameBasePlat, nil, mapstr.MapStr{}, &clouds, "create_time", 0, 0)
	if nil != err && !db.IsNotFoundError(err) {
		return err
	}

	flag := "updateflag"
	existDefault := false
	expects := map[string]map[string]interface{}{}
	for _, cloud := range clouds {
		if cloud[common.BKCloudNameField] == "default area" {
			cloud[common.BKCloudIDField] = 0
			existDefault = true
		}
		cloud[flag] = true
		expects[fmt.Sprintf("%v:%v", cloud[common.BKOwnerIDField], cloud[common.BKCloudNameField])] = cloud
	}

	if !existDefault {
		expects["0:"+"default area"] = map[string]interface{}{
			common.BKCloudNameField: "default area",
			common.BKOwnerIDField:   common.BKDefaultOwnerID,
			common.BKCloudIDField:   common.BKDefaultDirSubArea,
			common.CreateTimeField:  time.Now(),
			common.LastTimeField:    time.Now(),
			flag:                    true,
		}
	}

	for _, expect := range expects {
		if err = db.Table(common.BKTableNameBasePlat).Insert(ctx, expect); err != nil {
			return err
		}
	}

	if err = db.Table(common.BKTableNameBasePlat).Delete(ctx, mapstr.MapStr{
		flag: map[string]interface{}{
			common.BKDBNE: true,
		},
	}); err != nil {
		return err
	}

	if err = db.Table(common.BKTableNameBasePlat).DropColumn(ctx, flag); err != nil {
		return err
	}

	return nil
}
