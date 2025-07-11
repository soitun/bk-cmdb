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

package event

import (
	"time"

	"configcenter/src/common/metrics"
	"configcenter/src/storage/stream/types"
	"github.com/prometheus/client_golang/prometheus"
)

// InitialMetrics TODO
func InitialMetrics(collection string, subSys string) *EventMetrics {
	labels := prometheus.Labels{"collection": collection}

	m := new(EventMetrics)
	m.totalEventCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   subSys,
		Name:        "total_event_count",
		Help:        "the total event count which we handled with this resources",
		ConstLabels: labels,
	}, []string{"action"})
	metrics.Register().MustRegister(m.totalEventCount)

	m.lastEventTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   subSys,
		Name:        "last_event_unix_time_seconds",
		Help:        "records the time that event occurs at unix time seconds",
		ConstLabels: labels,
	}, []string{})
	metrics.Register().MustRegister(m.lastEventTime)

	m.eventLagDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   subSys,
		Name:        "event_lag_seconds",
		Help:        "the lags(seconds) of the event between it occurs and we received it",
		ConstLabels: labels,
		Buckets:     []float64{0.02, 0.04, 0.06, 0.08, 0.1, 0.3, 0.5, 0.7, 1, 5, 10, 20, 30, 60, 120},
	}, []string{"action"})
	metrics.Register().MustRegister(m.eventLagDurations)

	m.lastEventLagDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   subSys,
		Name:        "last_event_lag_seconds",
		Help:        "record the last event's lag duration in seconds",
		ConstLabels: labels,
	}, []string{})
	metrics.Register().MustRegister(m.lastEventLagDuration)

	m.cycleDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   metrics.Namespace,
		Subsystem:   subSys,
		Name:        "event_cycle_seconds",
		Help:        "the total duration(seconds) of each event being handled",
		ConstLabels: labels,
		Buckets:     []float64{0.02, 0.04, 0.06, 0.08, 0.1, 0.3, 0.5, 0.7, 1, 5, 10, 15, 20, 40, 60, 120},
	}, []string{})
	metrics.Register().MustRegister(m.cycleDurations)

	m.totalErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: subSys,
		Name:      "total_error_count",
		Help: "the total error event count which we handled with this resources, " +
			"including invalid event and re-watch operations",
		ConstLabels: labels,
	}, []string{"error_type"})
	metrics.Register().MustRegister(m.totalErrorCount)

	return m
}

// EventMetrics TODO
type EventMetrics struct {
	// record the total event count be handled.
	totalEventCount *prometheus.CounterVec

	// last event time records when the current last event occurs and when we received with unix time seconds.
	// we can use this metric to know which time point we have handled the event for now.
	lastEventTime *prometheus.GaugeVec

	// record the all event lag which is the duration time difference
	// between the event is occur and the time we received the with watch.
	// unit is seconds.
	eventLagDurations *prometheus.HistogramVec

	// record the last event's lag duration
	lastEventLagDuration *prometheus.GaugeVec

	// record the cost time of every cycle of handling a event chain.
	// unit is seconds
	cycleDurations *prometheus.HistogramVec

	// record the total errors when we watch the event.
	// which contains the number of as follows:
	// 1. watch failed count
	// 2. invalid event count
	totalErrorCount *prometheus.CounterVec
}

// CollectBasic collect the basic event's metrics
func (em *EventMetrics) CollectBasic(e *types.Event) {
	// increase event's total count with operate type
	em.totalEventCount.With(prometheus.Labels{"action": string(e.OperationType)}).Inc()

	// the time when the event is really happens.
	at := time.Unix(int64(e.ClusterTime.Sec), int64(e.ClusterTime.Nano))

	// unix time in seconds
	em.lastEventTime.With(prometheus.Labels{}).Set(float64(at.Unix()))

	// calculate event lags, in seconds
	lags := time.Since(at).Seconds()

	// set last lag duration
	em.lastEventLagDuration.With(prometheus.Labels{}).Set(lags)

	// add to lags durations
	em.eventLagDurations.With(prometheus.Labels{"action": string(e.OperationType)}).Observe(lags)
}

// CollectCycleDuration TODO
// the total duration(seconds) of each event being handled
func (em *EventMetrics) CollectCycleDuration(d time.Duration) {
	em.cycleDurations.With(prometheus.Labels{}).Observe(d.Seconds())
}

// CollectRetryError TODO
// collect retry operation for any reason
func (em *EventMetrics) CollectRetryError() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "retry"}).Inc()
}

// CollectRedisError TODO
// collect redis operation related errors
func (em *EventMetrics) CollectRedisError() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "redis_command"}).Inc()
}

// CollectMongoError TODO
// collect mongodb related errors, such as get info from table cc_DelArchive
func (em *EventMetrics) CollectMongoError() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "mongo_command"}).Inc()
}

// CollectConflict collect event conflict count
func (em *EventMetrics) CollectConflict() {
	em.totalErrorCount.With(prometheus.Labels{"error_type": "conflict"}).Inc()
}
