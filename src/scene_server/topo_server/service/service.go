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

// Package service TODO
package service

import (
	"encoding/json"

	"configcenter/src/ac/extensions"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/logics"
	"configcenter/src/thirdparty/elasticsearch"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// Service TODO
type Service struct {
	Engine      *backbone.Engine
	Logics      logics.Logics
	Config      options.Config
	AuthManager *extensions.AuthManager
	Es          *elasticsearch.EsSrv
	Error       errors.CCErrorIf
	Language    language.CCLanguageIf
}

// WebService the web service
func (s *Service) WebService() *restful.Container {
	errors.SetGlobalCCError(s.Error)
	getErrFunc := func() errors.CCErrorIf {
		return s.Error
	}

	api := new(restful.WebService)
	api.Path("/topo/v3/").Filter(s.Engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).
		Produces(restful.MIME_JSON)

	// init service actions
	s.initService(api)

	container := restful.NewContainer().Add(api)

	opentelemetry.AddOtlpFilter(container)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	container.Add(commonAPI)

	return container
}

// ModelType is model type
// bk_biz_id == 0 : public model
// bk_biz_id > 0 : private model
type ModelType struct {
	BizID int64 `json:"bk_biz_id"`
}

// MapStrWithModelBizID TODO
type MapStrWithModelBizID struct {
	ModelBizID int64
	Data       mapstr.MapStr
}

// UnmarshalJSON TODO
func (m *MapStrWithModelBizID) UnmarshalJSON(data []byte) error {
	modelType := new(ModelType)
	if err := json.Unmarshal(data, modelType); err != nil {
		return err
	}
	m.ModelBizID = modelType.BizID
	if err := json.Unmarshal(data, &m.Data); err != nil {
		return err
	}

	return nil
}
