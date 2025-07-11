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

package v3v0v8

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
)

const (
	tableNameSubscription = "cc_Subscription"
	subscriptionIDField   = "subscription_id"
)

func createTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	for tableName, indexes := range tables {
		exists, err := db.HasTable(ctx, tableName)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(ctx, tableName); err != nil && !mgo.IsDup(err) {
				return err
			}
		}
		for index := range indexes {
			if err = db.Table(tableName).CreateIndex(ctx, indexes[index]); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
	}
	return nil
}

var tables = map[string][]types.Index{
	common.BKTableNameBaseApp: {
		types.Index{Name: "", Keys: bson.D{{common.BKAppIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKAppNameField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKDefaultField, 1}}, Background: true},
	},

	common.BKTableNameBaseHost: {
		types.Index{Name: "", Keys: bson.D{{common.BKHostIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKHostNameField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKHostInnerIPField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKHostOuterIPField, 1}}, Background: true},
	},
	common.BKTableNameBaseModule: {
		types.Index{Name: "", Keys: bson.D{{common.BKModuleIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKModuleNameField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKDefaultField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKAppIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKSetIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKParentIDField, 1}}, Background: true},
	},
	common.BKTableNameModuleHostConfig: {
		types.Index{Name: "", Keys: bson.D{{common.BKAppIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKHostIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKModuleIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKSetIDField, 1}}, Background: true},
	},
	common.BKTableNameObjAsst: {
		types.Index{Name: "", Keys: bson.D{{common.BKObjIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKAsstObjIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
	},
	common.BKTableNameObjAttDes: {
		types.Index{Name: "", Keys: bson.D{{common.BKObjIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKFieldID, 1}}, Background: true},
	},
	common.BKTableNameObjClassification: {
		types.Index{Name: "", Keys: bson.D{{common.BKClassificationIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKClassificationNameField, 1}}, Background: true},
	},
	common.BKTableNameObjDes: {
		types.Index{Name: "", Keys: bson.D{{common.BKObjIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKClassificationIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKObjNameField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
	},
	common.BKTableNameBaseInst: {
		types.Index{Name: "", Keys: bson.D{{common.BKObjIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKInstIDField, 1}}, Background: true},
	},
	common.BKTableNameAuditLog: {
		{Name: "index_bk_supplier_account", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
		{Name: "index_audit_type", Keys: bson.D{{common.BKAuditTypeField, 1}}, Background: true},
		{Name: "index_action", Keys: bson.D{{common.BKActionField, 1}}, Background: true},
	},
	common.BKTableNameBasePlat: {
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
	},
	"cc_Proc2Module": {
		types.Index{Name: "", Keys: bson.D{{common.BKAppIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKProcIDField, 1}}, Background: true},
	},
	common.BKTableNameBaseProcess: {
		types.Index{Name: "", Keys: bson.D{{common.BKProcIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKAppIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
	},
	common.BKTableNamePropertyGroup: {
		types.Index{Name: "", Keys: bson.D{{common.BKObjIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKPropertyGroupIDField, 1}}, Background: true},
	},
	common.BKTableNameBaseSet: {
		types.Index{Name: "", Keys: bson.D{{common.BKSetIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKParentIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKAppIDField, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BkSupplierAccount, 1}}, Background: true},
		types.Index{Name: "", Keys: bson.D{{common.BKSetNameField, 1}}, Background: true},
	},
	tableNameSubscription: {
		types.Index{Name: "", Keys: bson.D{{subscriptionIDField, 1}}, Background: true},
	},
	common.BKTableNameTopoGraphics: {
		types.Index{Name: "", Keys: bson.D{{"scope_type", 1}, {"scope_id", 1}, {"node_type", 1},
			{common.BKObjIDField, 1}, {common.BKInstIDField, 1}}, Background: true, Unique: true},
	},
	common.BKTableNameInstAsst: {
		types.Index{Name: "", Keys: bson.D{{common.BKObjIDField, 1}, {common.BKInstIDField, 1}},
			Background: true},
	},

	common.BKTableNameHistory:      {},
	common.BKTableNameHostFavorite: {},
	common.BKTableNameUserAPI:      {},
	common.BKTableNameUserCustom:   {},
	common.BKTableNameIDgenerator:  {},
	common.BKTableNameSystem:       {},
	common.BKTableNameDelArchive:   {},
}
