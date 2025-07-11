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

// UpdateHostCloudAreaField TODO
func (s *coreService) UpdateHostCloudAreaField(ctx *rest.Contexts) {
	input := metadata.UpdateHostCloudAreaFieldOption{}
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().UpdateHostCloudAreaField(ctx.Kit, input)
	if err != nil {
		blog.Errorf("UpdateHostCloudAreaField failed, call core operation failed, input: %+v, err: %v, rid: %v", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// FindCloudAreaHostCount TODO
func (s *coreService) FindCloudAreaHostCount(ctx *rest.Contexts) {
	input := metadata.CloudAreaHostCount{}
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	res, err := s.core.HostOperation().FindCloudAreaHostCount(ctx.Kit, input)
	if err != nil {
		blog.Errorf("UpdateHostCloudAreaField failed, call core operation failed, input: %+v, err: %v, rid: %v", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}
