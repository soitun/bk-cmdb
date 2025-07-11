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

package local

import (
	"testing"

	"configcenter/src/common/metadata"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDecodeSubMap(t *testing.T) {

	ts := metadata.Now()
	attr := metadata.Attribute{
		BizID:      111,
		LastTime:   &ts,
		CreateTime: &ts,
	}

	bsonAttr, err := bson.Marshal(attr)
	require.NoError(t, err)

	newAttr := &metadata.Attribute{}
	err = bson.Unmarshal(bsonAttr, newAttr)
	require.NoError(t, err)

	require.Equal(t, attr.BizID, newAttr.BizID)

}
