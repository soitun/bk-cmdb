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

package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
)

// PubSub is the interface for redis Pub/Sub commands
type PubSub interface {
	Channel() <-chan *redis.Message
	ChannelSize(size int) <-chan *redis.Message
	ChannelWithSubscriptions(size int) <-chan interface{}
	Close() error
	PSubscribe(patterns ...string) error
	PUnsubscribe(patterns ...string) error
	Ping(payload ...string) error
	Receive() (interface{}, error)
	ReceiveMessage() (*redis.Message, error)
	ReceiveTimeout(timeout time.Duration) (interface{}, error)
	String() string
	Subscribe(channels ...string) error
	Unsubscribe(channels ...string) error
}
