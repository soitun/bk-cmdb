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

// Package apigw defines api-gateway client
package apigw

import (
	"configcenter/src/common/blog"
	"configcenter/src/thirdparty/apigw"
	"configcenter/src/thirdparty/apigw/apigwutil"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	client apigw.ClientSet
)

// Client returns api-gateway client set
func Client() apigw.ClientSet {
	return client
}

// Init the api-gateway client set
func Init(path string, metric prometheus.Registerer, neededClients []apigw.ClientType) error {
	config, err := apigwutil.ParseApiGWConfig(path)
	if err != nil {
		blog.Errorf("parse %s api gateway config failed, err: %v", path, err)
		return err
	}

	client, err = apigw.NewClientSet(config, metric, neededClients)
	if err != nil {
		blog.Errorf("init %s api gateway client failed, err: %v", path, err)
		return err
	}

	return nil
}
