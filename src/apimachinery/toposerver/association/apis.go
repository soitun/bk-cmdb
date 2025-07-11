// Package association TODO
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
package association

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

// AssociationInterface TODO
type AssociationInterface interface {
	SearchType(ctx context.Context, h http.Header,
		request *metadata.SearchAssociationTypeRequest) (resp *metadata.SearchAssociationTypeResult, err error)
	CreateType(ctx context.Context, h http.Header,
		request *metadata.AssociationKind) (resp *metadata.CreateAssociationTypeResult, err error)
	UpdateType(ctx context.Context, h http.Header, asstTypeID int,
		request *metadata.UpdateAssociationTypeRequest) (resp *metadata.UpdateAssociationTypeResult, err error)
	DeleteType(ctx context.Context, h http.Header, asstTypeID int) (resp *metadata.DeleteAssociationTypeResult,
		err error)
	SearchObject(ctx context.Context, h http.Header,
		request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error)
	CreateObject(ctx context.Context, h http.Header,
		request *metadata.Association) (resp *metadata.CreateAssociationObjectResult, err error)
	UpdateObject(ctx context.Context, h http.Header, asstID int,
		request *metadata.UpdateAssociationObjectRequest) (resp *metadata.UpdateAssociationObjectResult, err error)
	DeleteObject(ctx context.Context, h http.Header, asstID int) (resp *metadata.DeleteAssociationObjectResult,
		err error)
	SearchInst(ctx context.Context, h http.Header,
		request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error)
	SearchAssociationRelatedInst(ctx context.Context, h http.Header,
		request *metadata.SearchAssociationRelatedInstRequest) (resp *metadata.SearchAssociationInstResult, err error)
	CreateInst(ctx context.Context, h http.Header,
		request *metadata.CreateAssociationInstRequest) (resp *metadata.CreateAssociationInstResult, err error)
	CreateManyInstAssociation(ctx context.Context, header http.Header,
		request *metadata.CreateManyInstAsstRequest) (*metadata.CreateManyInstAsstResult, error)
	DeleteInst(ctx context.Context, h http.Header, objID string,
		assoID int64) (resp *metadata.DeleteAssociationInstResult, err error)
	DeleteInstBatch(ctx context.Context, h http.Header,
		assoIDs *metadata.DeleteAssociationInstBatchRequest) (resp *metadata.DeleteAssociationInstBatchResult,
		err error)
	SearchObjectAssoWithAssoKindList(ctx context.Context, h http.Header,
		assoKindIDs metadata.AssociationKindIDs) (resp *metadata.ListAssociationsWithAssociationKindResult, err error)

	SearchInstAssocAndInstDetail(ctx context.Context, header http.Header, objID string,
		input *metadata.InstAndAssocRequest) (*metadata.InstAndAssocDetailResult, error)

	// SearchInstanceAssociations searches object instance associations.
	SearchInstanceAssociations(ctx context.Context, header http.Header,
		objID string, input *metadata.CommonSearchFilter) (*metadata.Response, error)

	// CountInstanceAssociations counts object instance associations num.
	CountInstanceAssociations(ctx context.Context, header http.Header,
		objID string, input *metadata.CommonCountFilter) (*metadata.Response, error)
}

// NewAssociationInterface TODO
func NewAssociationInterface(client rest.ClientInterface) AssociationInterface {
	return &Association{client: client}
}

// Association TODO
type Association struct {
	client rest.ClientInterface
}
