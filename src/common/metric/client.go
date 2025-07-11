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

package metric

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

var metricController *MetricController

func newMetricController(conf Config, healthFunc HealthFunc, collectors ...*Collector) []Action {
	metricController = new(MetricController)
	meta := MetaData{
		Module:        conf.ModuleName,
		ServerAddress: conf.ServerAddress,
		Labels:        conf.Labels,
	}

	// set default golang metric.
	collectors = append(collectors, newGoMetricCollector())

	metricController.MetaData = &meta
	metricController.Collectors = make(map[CollectorName]CollectInter)
	for _, c := range collectors {
		metricController.Collectors[c.Name] = c.Collector
	}

	metricHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metric, err := metricController.PackMetrics()
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			info := fmt.Sprintf("get metrics failed. err: %v", err)
			w.Write([]byte(info))
			return
		}
		w.Write(*metric)
	})

	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h := healthFunc()
		info := HealthInfo{
			Module:     conf.ModuleName,
			Address:    conf.ServerAddress,
			HealthMeta: h,
			AtTime:     metadata.Now(),
		}

		rsp := HealthResponse{
			Code:    common.CCSuccess,
			Data:    info,
			OK:      h.IsHealthy,
			Result:  h.IsHealthy,
			Message: h.Message,
		}
		rsp.SetCommonResponse()
		js, err := json.Marshal(rsp)
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			info := fmt.Sprintf("get health info failed. err: %v", err)
			w.Write([]byte(info))
			return
		}
		w.Write(js)
	})

	actions := []Action{
		{Method: "GET", Path: "/metrics", HandlerFunc: metricHandler},
		{Method: "GET", Path: "/healthz", HandlerFunc: healthHandler},
	}

	return actions
}

// MetricController TODO
type MetricController struct {
	MetaData   *MetaData
	Collectors map[CollectorName]CollectInter
}

// PackMetrics TODO
func (mc *MetricController) PackMetrics() (*[]byte, error) {
	mf := MetricFamily{
		MetaData:     mc.MetaData,
		MetricBundle: make(map[CollectorName][]*Metric),
	}

	for name, collector := range mc.Collectors {
		mf.MetricBundle[name] = make([]*Metric, 0)
		done := make(chan struct{}, 0)
		go func(c CollectInter) {
			for _, mc := range c.Collect() {
				metric, err := newMetric(mc)
				if nil != err {
					blog.Errorf("new metric failed. err: %v", err)
					continue
				}
				mf.MetricBundle[name] = append(mf.MetricBundle[name], metric)
			}
			done <- struct{}{}
		}(collector)

		select {
		case <-time.After(time.Duration(10 * time.Second)):
			blog.Errorf("get metric bundle: %s timeout, skip it.", name)
			continue
		case <-done:
			close(done)
		}
	}

	mf.ReportTimeMs = time.Now().Unix()
	js, err := json.Marshal(mf)
	if nil != err {
		return nil, err
	}
	return &js, nil
}
