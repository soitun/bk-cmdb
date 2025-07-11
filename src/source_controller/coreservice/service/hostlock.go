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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// LockHost TODO
func (s *coreService) LockHost(ctx *rest.Contexts) {
	input := new(metadata.HostLockRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().LockHost(ctx.Kit, input)
	if err != nil {
		blog.Errorf("LockHost failed, lock host handle failed, err: %+v, input:%+v, rid:%s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// UnlockHost TODO
func (s *coreService) UnlockHost(ctx *rest.Contexts) {
	input := new(metadata.HostLockRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().UnlockHost(ctx.Kit, input)
	if err != nil {
		blog.Errorf("UnlockHost failed, unlock host handle failed, err: %s, input:%+v, rid:%s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// QueryLockHost TODO
func (s *coreService) QueryLockHost(ctx *rest.Contexts) {
	input := new(metadata.QueryHostLockRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	hostLockArr, err := s.core.HostOperation().QueryHostLock(ctx.Kit, input)
	if err != nil {
		blog.Errorf("QueryLockHost failed, query host handle failed, err: %s, input:%+v, rid: %s", err.Error(), input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	result := metadata.HostLockQueryResponse{}
	result.Data.Info = hostLockArr
	result.Data.Count = int64(len(hostLockArr))
	ctx.RespEntity(result.Data)
}
