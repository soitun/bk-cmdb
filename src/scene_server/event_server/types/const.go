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

package types

const (
	// NilStr is special NIL string.
	NilStr = "nil"

	// MetricsNamespacePrefix is prefix of metrics namespace.
	MetricsNamespacePrefix = "cmdb_eventserver"
)

// ApiVersion gse api version
type ApiVersion string

const (
	// V1 use gse1.0 thrift api
	V1 ApiVersion = "v1"
	// V2 use gse2.0 apiGW api
	V2 ApiVersion = "v2"
)
