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

package inst

import (
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

// CreateApp TODO
func (t *instanceClient) CreateApp(ctx context.Context, ownerID string, h http.Header,
	params map[string]interface{}) (resp *metadata.CreateInstResult, err error) {
	resp = new(metadata.CreateInstResult)
	subPath := "/app/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteApp TODO
func (t *instanceClient) DeleteApp(ctx context.Context, ownerID string, appID string,
	h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/app/%s/%s"

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateApp TODO
func (t *instanceClient) UpdateApp(ctx context.Context, ownerID string, appID string, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/app/%s/%s"
	err = t.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateAppDataStatus TODO
func (t *instanceClient) UpdateAppDataStatus(ctx context.Context, ownerID string, flag common.DataStatusFlag,
	appID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/app/%s/%s/%s"
	err = t.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, flag, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchApp TODO
func (t *instanceClient) SearchApp(ctx context.Context, ownerID string, h http.Header,
	s *params.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := "/app/search/%s"
	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetAppBasicInfo TODO
func (t *instanceClient) GetAppBasicInfo(ctx context.Context, h http.Header,
	bizID int64) (resp *metadata.AppBasicInfoResult, err error) {
	resp = new(metadata.AppBasicInfoResult)
	subPath := "/app/%d/basic_info"
	err = t.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetDefaultApp TODO
func (t *instanceClient) GetDefaultApp(ctx context.Context, ownerID string,
	h http.Header) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := "/app/default/%s/search"
	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// CreateDefaultApp TODO
func (t *instanceClient) CreateDefaultApp(ctx context.Context, ownerID string, h http.Header,
	data map[string]interface{}) (resp *metadata.CreateInstResult, err error) {
	resp = new(metadata.CreateInstResult)
	subPath := "/app/default/%s"
	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
