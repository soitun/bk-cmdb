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

package mapping

import (
	"fmt"

	"configcenter/pkg/cache/general"
	"configcenter/src/common/watch"
)

var cursorTypeMap = map[general.ResType]watch.CursorType{
	general.Host:             watch.Host,
	general.ModuleHostRel:    watch.ModuleHostRelation,
	general.Biz:              watch.Biz,
	general.Set:              watch.Set,
	general.Module:           watch.Module,
	general.Process:          watch.Process,
	general.ProcessRelation:  watch.ProcessInstanceRelation,
	general.BizSet:           watch.BizSet,
	general.Plat:             watch.Plat,
	general.Project:          watch.Project,
	general.ObjectInstance:   watch.ObjectBase,
	general.MainlineInstance: watch.MainlineInstance,
	general.InstAsst:         watch.InstAsst,
	general.KubeCluster:      watch.KubeCluster,
	general.KubeNode:         watch.KubeNode,
	general.KubeNamespace:    watch.KubeNamespace,
	general.KubeWorkload:     watch.KubeWorkload,
	general.KubePod:          watch.KubePod,
}

// GetCursorTypeByResType get event watch cursor type by resource type
func GetCursorTypeByResType(res general.ResType) (watch.CursorType, error) {
	typ, exists := cursorTypeMap[res]
	if !exists {
		return watch.UnknownType, fmt.Errorf("resource type %s has no matching cursor type", res)
	}

	return typ, nil
}
