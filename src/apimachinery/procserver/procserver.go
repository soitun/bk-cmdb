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

// Package procserver TODO
package procserver

import (
	"fmt"

	"configcenter/src/apimachinery/procserver/openapi"
	"configcenter/src/apimachinery/procserver/process"
	"configcenter/src/apimachinery/procserver/service"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

// ProcServerClientInterface TODO
type ProcServerClientInterface interface {
	Process() process.ProcessClientInterface
	OpenAPI() openapi.OpenAPIClientInterface
	Service() service.ServiceClientInterface
}

// NewProcServerClientInterface TODO
func NewProcServerClientInterface(c *util.Capability, version string) ProcServerClientInterface {
	base := fmt.Sprintf("/process/%s", version)
	return &procServer{client: rest.NewRESTClient(c, base)}
}

type procServer struct {
	client rest.ClientInterface
}

// Process TODO
func (p *procServer) Process() process.ProcessClientInterface {
	return process.NewProcessClientInterface(p.client)
}

// OpenAPI TODO
func (p *procServer) OpenAPI() openapi.OpenAPIClientInterface {
	return openapi.NewOpenApiClientInterface(p.client)
}

// Service TODO
func (p *procServer) Service() service.ServiceClientInterface {
	return service.NewServiceClientInterface(p.client)
}
