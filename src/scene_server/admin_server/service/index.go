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

package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/storage/dal/mongo/local"

	"github.com/emicklei/go-restful/v3"
)

// BackgroundTask TODO
func (s *Service) BackgroundTask(options options.Config) error {

	mongoConf, err := cc.Mongo("mongodb")
	if err != nil {
		return err
	}

	// db 语句的执行时间设置为never timeout
	mongoConf.SocketTimeout = 0
	db, err := local.NewMgo(mongoConf.GetMongoConf(), time.Minute)
	if err != nil {
		return fmt.Errorf("connect mongo server failed %s", err.Error())
	}

	logics.DBSync(s.Engine, db, options)

	return nil
}

// RunSyncDBIndex TODO
func (s *Service) RunSyncDBIndex(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	rid := httpheader.GetRid(rHeader)
	ctx := context.WithValue(context.Background(), common.ContextRequestIDField, rid)

	if err := logics.RunSyncDBIndex(ctx, s.Engine); err != nil {
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(""))

	return
}
