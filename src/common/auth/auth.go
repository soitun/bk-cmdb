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

// Package auth TODO
package auth

import (
	"strconv"
	"sync"

	"configcenter/src/common/blog"
)

// EnableAuth TODO
var EnableAuth = "true"
var enableAuth = true

// EnableAuthFlag TODO
var EnableAuthFlag *authValue
var once = sync.Once{}

type authValue struct{}

// String 用于打印
func (a *authValue) String() string {
	return strconv.FormatBool(enableAuth)
}

// Set TODO
func (a *authValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	setEnableAuth(v)
	return nil
}

// Type TODO
func (a *authValue) Type() string {
	return "bool"
}

func init() {
	var err error
	enableAuth, err = strconv.ParseBool(EnableAuth)
	if err != nil {
		blog.Fatalf("[auth] enableAuth %s configuration invalid", EnableAuth)
	}
}

// setEnableAuth is the default handler which match the --enable-auth flag
func setEnableAuth(enable bool) {
	once.Do(func() {
		enableAuth = enable
	})
}

// EnableAuthorize TODO
func EnableAuthorize() bool {
	return enableAuth
}
