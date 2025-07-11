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
	"context"
	"net/http"

	"configcenter/src/ac"
	"configcenter/src/ac/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/event_server/sync/hostidentifier"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// Service impls main logics as service.
type Service struct {
	ctx    context.Context
	engine *backbone.Engine

	db          dal.RDB
	cache       redis.Client
	authorizer  ac.AuthorizeInterface
	AuthManager *extensions.AuthManager

	// SyncData is sync host identifier operator
	SyncData *hostidentifier.HostIdentifier
}

// NewService creates a new Service object.
func NewService(ctx context.Context, engine *backbone.Engine) *Service {
	return &Service{ctx: ctx, engine: engine}
}

// SetDB setups database.
func (s *Service) SetDB(db dal.RDB) {
	s.db = db
}

// SetCache setups cc main redis.
func (s *Service) SetCache(db redis.Client) {
	s.cache = db
}

// SetAuthorizer TODO
func (s *Service) SetAuthorizer(authorizer ac.AuthorizeInterface) {
	s.authorizer = authorizer
}

// WebService setups a new restful web service.
func (s *Service) WebService() *restful.Container {
	container := restful.NewContainer()

	opentelemetry.AddOtlpFilter(container)

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.engine.CCErr
	}

	api.Path("/event/v3")
	api.Filter(s.engine.Metric().RestfulMiddleWare)
	api.Filter(rdapi.AllGlobalFilter(getErrFunc))
	api.Produces(restful.MIME_JSON)

	s.initService(api)

	container.Add(api)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	container.Add(commonAPI)

	return container
}

func (s *Service) initService(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/watch/resource/{resource}", Handler: s.WatchEvent})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/sync/host_identifier", Handler: s.SyncHostIdentifier})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/push/host_identifier", Handler: s.PushHostIdentifier})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/host_identifier_push_result",
		Handler: s.GetHostIdentifierPushResult})

	utility.AddToRestfulWebService(web)

}

// Healthz is a HTTP restful interface for health check.
func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	// metadata.
	meta := metric.HealthMeta{IsHealthy: true}

	// zookeeper health status info.
	zkItem := metric.HealthItem{
		IsHealthy: true,
		Name:      types.CCFunctionalityServicediscover,
	}
	if err := s.engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// mongodb health status info.
	meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityMongo, s.db.Ping()))

	// cc main redis health status info.
	meta.Items = append(meta.Items,
		metric.NewHealthItem(types.CCFunctionalityRedis, s.cache.Ping(context.Background()).Err()))

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "event server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_EVENTSERVER,
		HealthMeta: meta,
		AtTime:     metadata.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	answer.SetCommonResponse()

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}
