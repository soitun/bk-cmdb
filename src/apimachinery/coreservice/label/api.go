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

package label

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/selector"
	"configcenter/src/common/util"
)

// AddLabel TODO
func (l *label) AddLabel(ctx context.Context, h http.Header, tableName string,
	option selector.LabelAddOption) errors.CCErrorCoder {
	rid := util.ExtractRequestIDFromContext(ctx)
	ret := new(metadata.BaseResp)
	subPath := "/createmany/labels"

	body := selector.LabelAddRequest{
		Option:    option,
		TableName: tableName,
	}

	err := l.client.Post().
		WithContext(ctx).
		Body(body).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("AddLabel failed, http request failed, err: %+v, rid: %s", err, rid)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// RemoveLabel TODO
func (l *label) RemoveLabel(ctx context.Context, h http.Header, tableName string,
	option selector.LabelRemoveOption) errors.CCErrorCoder {
	rid := util.ExtractRequestIDFromContext(ctx)
	ret := new(metadata.BaseResp)
	subPath := "/deletemany/labels"

	body := selector.LabelRemoveRequest{
		Option:    option,
		TableName: tableName,
	}
	err := l.client.Delete().
		Body(body).
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("RemoveLabel failed, http request failed, err: %+v, rid: %s", err, rid)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// UpdateLabel update service instance tag request.
func (l *label) UpdateLabel(ctx context.Context, h http.Header, tableName string,
	option *selector.LabelUpdateOption) errors.CCErrorCoder {
	ret := new(metadata.BaseResp)
	subPath := "/updatemany/labels"

	body := selector.LabelUpdateRequest{
		Option:    option,
		TableName: tableName,
	}
	err := l.client.Post().
		Body(body).
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}
