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

package transaction

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

// NewTxn TODO
func NewTxn(client rest.ClientInterface) Interface {
	return &txn{client: client}
}

type txn struct {
	client rest.ClientInterface
}

// Interface TODO
// Transaction interface
type Interface interface {
	// NewTransaction 开启新事务
	NewTransaction(h http.Header, opts ...metadata.TxnOption) (Transaction, error)
	// AutoRunTxn is a transaction wrapper. it will automatically commit or abort the
	// transaction depend on the f(), if f() returns with an error, then abort the
	// transaction, otherwise, it will commit the transaction.
	AutoRunTxn(ctx context.Context, h http.Header, run func() error, opts ...metadata.TxnOption) error
}
