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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	ccjson "configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"github.com/emicklei/go-restful/v3"
)

// Get TODO
func (s *service) Get(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

// Put TODO
func (s *service) Put(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

// Post TODO
func (s *service) Post(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

// Delete TODO
func (s *service) Delete(req *restful.Request, resp *restful.Response) {
	s.Do(req, resp)
}

const maxToleranceLatencyTime = 5 * time.Second

// Do TODO
func (s *service) Do(req *restful.Request, resp *restful.Response) {

	rid := httpheader.GetRid(req.Request.Header)
	start := time.Now()
	url := req.Request.URL.Scheme + "://" + req.Request.URL.Host + req.Request.RequestURI
	proxyReq, err := http.NewRequestWithContext(req.Request.Context(), req.Request.Method, url, req.Request.Body)
	if err != nil {
		blog.Errorf("new proxy request[%s] failed, err: %v, rid: %s", url, err, rid)
		s.RespError(req, resp, http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("proxy request failed, %s", err.Error()),
			ErrCode: common.CCErrProxyRequestFailed,
			Data:    nil,
		})

		return
	}

	for k, v := range req.Request.Header {
		if len(v) > 0 {
			proxyReq.Header.Set(k, v[0])
		}
	}

	response, err := s.client.Do(proxyReq)
	if err != nil {
		if time.Since(start) >= maxToleranceLatencyTime {
			if !strings.Contains(req.Request.RequestURI, "/watch/resource/") {
				// except resource watch api
				blog.Warnf("request exceeded max latency time, %s, %s, cost: %d ms, rid: %s", req.Request.Method, url,
					time.Since(start)/time.Millisecond, rid)
			}
		}

		blog.Errorf("*failed to do request[%s url: %s], user: %s, app code: %s, err: %v, rid: %s", req.Request.Method,
			url, httpheader.GetUser(req.Request.Header), httpheader.GetAppCode(req.Request.Header),
			err, rid)

		// send alarm when http request timeout, to monitor api server request
		if strings.Contains(err.Error(), "timeout awaiting response headers") {
			monitor.Collect(&meta.Alarm{
				RequestID: rid,
				Type:      meta.HttpFatalError,
				Detail: fmt.Sprintf("request timeout, user: %s, app code: %s, err: %v, %s, %s, rid: %s, cost: %d ms",
					httpheader.GetUser(req.Request.Header),
					httpheader.GetAppCode(req.Request.Header), err, req.Request.Method, url, rid,
					time.Since(start)/time.Millisecond),
				Dimension: map[string]string{"error_type": "request timeout"},
				Module:    common.GetIdentification(),
			})
		}

		s.RespError(req, resp, http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("proxy request failed, %s", err.Error()),
			ErrCode: common.CCErrProxyRequestFailed,
			Data:    nil,
		})
		return
	}

	if time.Since(start) >= maxToleranceLatencyTime {
		if !strings.Contains(req.Request.RequestURI, "/watch/resource/") {
			// except resource watch api
			blog.Warnf("request exceeded max latency time, %s, %s, cost: %d ms, rid: %s", req.Request.Method, url,
				time.Since(start)/time.Millisecond, rid)
		}
	}

	for k, v := range response.Header {
		if len(v) > 0 {
			resp.Header().Set(k, v[0])
		}
	}

	parseResponse(req, resp, response.Body, rid)

	blog.V(4).Infof("cost: %dms, action: %s, status code: %d, user: %s, app code: %s, url: %s, rid: %s",
		time.Since(start).Nanoseconds()/int64(time.Millisecond), req.Request.Method, response.StatusCode,
		httpheader.GetUser(req.Request.Header), httpheader.GetAppCode(req.Request.Header), url,
		httpheader.GetRid(req.Request.Header),
	)
	return
}

func parseResponse(req *restful.Request, resp *restful.Response, body io.ReadCloser, rid string) {
	// compatible for esb and old ui response
	// TODO remove this logics and change cc response format when esb is not supported
	header := req.Request.Header
	if httpheader.GetBkJWT(header) == "" || httpheader.IsReqFromWeb(header) {
		if _, err := io.Copy(resp, body); err != nil {
			body.Close()
			blog.Errorf("response request[url: %s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
			return
		}
		body.Close()
		return
	}

	// parse api gateway response, change the format to esb style
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		blog.Errorf("read response body failed, err: %v, rid: %s", err, rid)
		return
	}
	defer body.Close()

	keyMap := map[string]string{
		common.HTTPBKAPIErrorCode:    common.BkAPIErrorCode,
		common.HTTPBKAPIErrorMessage: common.BkAPIErrorMessage,
	}
	convertedBody, err := ccjson.ReplaceJsonKey(bodyBytes, keyMap)
	if err != nil {
		blog.Errorf("replace resp(%s) key failed, err: %v, rid: %s", string(bodyBytes), err, rid)
		return
	}

	resp.Header().Set("Content-Length", strconv.Itoa(len(convertedBody)))

	_, err = resp.Write(convertedBody)
	if err != nil {
		blog.Errorf("write api gateway response failed, err: %v, rid: %s", err, rid)
		return
	}
}
