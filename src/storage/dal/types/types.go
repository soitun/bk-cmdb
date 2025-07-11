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

// Package types TODO
package types

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// Errors defines
var (
	ErrTransactionStated   = errors.New("transaction already started")
	ErrTransactionNotFound = errors.New("not in transaction environment")
	ErrDocumentNotFound    = errors.New("document not found")
	ErrDuplicated          = errors.New("duplicated")
	ErrSessionNotStarted   = errors.New("session is not started")

	UpdateOpAddToSet = "addToSet"
	UpdateOpPull     = "pull"
)

// Filter condition alias name
type Filter interface{}

// Table TODO
type Table interface {
	// Find 查询多个并反序列化到 Result
	Find(filter Filter, opts ...*FindOpts) Find
	// AggregateOne 聚合查询
	AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error
	AggregateAll(ctx context.Context, pipeline interface{}, result interface{}, opts ...*AggregateOpts) error
	// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
	Insert(ctx context.Context, docs interface{}) error
	// Update 更新数据
	Update(ctx context.Context, filter Filter, doc interface{}) error
	// Upsert TODO
	// update or insert data
	Upsert(ctx context.Context, filter Filter, doc interface{}) error
	// UpdateMultiModel  data based on operators.
	UpdateMultiModel(ctx context.Context, filter Filter, updateModel ...ModeUpdate) error

	// Delete 删除数据
	Delete(ctx context.Context, filter Filter) error

	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, index Index) error
	// BatchCreateIndexes 批量创建索引
	BatchCreateIndexes(ctx context.Context, index []Index) error

	// DropIndex 移除索引
	DropIndex(ctx context.Context, indexName string) error
	// Indexes 查询索引
	Indexes(ctx context.Context) ([]Index, error)

	// AddColumn 添加字段
	AddColumn(ctx context.Context, column string, value interface{}) error
	// RenameColumn 重命名字段
	RenameColumn(ctx context.Context, filter Filter, oldName, newColumn string) error
	// DropColumn 移除字段
	DropColumn(ctx context.Context, field string) error
	// DropColumns 根据条件移除字段
	DropColumns(ctx context.Context, filter Filter, fields []string) error

	// DropDocsColumn remove a column by the name for doc use filter
	DropDocsColumn(ctx context.Context, field string, filter Filter) error

	// Distinct Finds the distinct values for a specified field across a single collection or view and returns the results in an
	// field the field for which to return distinct values.
	// filter query that specifies the documents from which to retrieve the distinct values.
	Distinct(ctx context.Context, field string, filter Filter) ([]interface{}, error)

	// DeleteMany delete document, return number of documents that were deleted.
	DeleteMany(ctx context.Context, filter Filter) (uint64, error)
	// UpdateMany update document, return number of documents that were modified.
	UpdateMany(ctx context.Context, filter Filter, doc interface{}) (uint64, error)
}

// Find find operation interface
type Find interface {
	// Fields 设置查询字段
	Fields(fields ...string) Find
	// Sort 设置查询排序
	Sort(sort string) Find
	// Start 设置限制查询上标
	Start(start uint64) Find
	// Limit 设置查询数量
	Limit(limit uint64) Find
	// All 查询多个
	All(ctx context.Context, result interface{}) error
	// One 查询单个
	One(ctx context.Context, result interface{}) error
	// Count 统计数量(非事务)
	Count(ctx context.Context) (uint64, error)
	// List 查询多个, start 等于0的时候，返回满足条件的行数
	List(ctx context.Context, result interface{}) (int64, error)

	Option(opts ...*FindOpts)
}

// ModeUpdate  根据不同的操作符去更新数据
type ModeUpdate struct {
	Op  string
	Doc interface{}
}

// Index define the DB index struct
type Index struct {
	Keys                    bson.D                 `json:"keys" bson:"key"`
	Name                    string                 `json:"name" bson:"name"`
	Unique                  bool                   `json:"unique" bson:"unique"`
	Background              bool                   `json:"background" bson:"background"`
	ExpireAfterSeconds      int32                  `json:"expire_after_seconds" bson:"expire_after_seconds,omitempty"`
	PartialFilterExpression map[string]interface{} `json:"partialFilterExpression" bson:"partialFilterExpression"`
}

// FindOpts TODO
type FindOpts struct {
	WithObjectID *bool
	WithCount    *bool
}

// NewFindOpts TODO
func NewFindOpts() *FindOpts {
	return &FindOpts{}
}

// SetWithObjectID TODO
func (f *FindOpts) SetWithObjectID(bl bool) *FindOpts {
	f.WithObjectID = &bl
	return f
}

// SetWithCount TODO
func (f *FindOpts) SetWithCount(bl bool) *FindOpts {
	f.WithCount = &bl
	return f
}

// AggregateOpts TODO
type AggregateOpts struct {
	AllowDiskUse *bool
}

// NewAggregateOpts TODO
func NewAggregateOpts() *AggregateOpts {
	return &AggregateOpts{}
}

// SetAllowDiskUse TODO
func (a *AggregateOpts) SetAllowDiskUse(bl bool) *AggregateOpts {
	a.AllowDiskUse = &bl
	return a
}
