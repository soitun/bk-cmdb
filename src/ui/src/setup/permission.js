/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import cursor from '@/directives/cursor'
import { IAM_ACTIONS, OPERATION } from '@/dictionary/iam-auth'
import { $error } from '@/magicbox'
import isEqual from 'lodash/isEqual'
import uniqBy from 'lodash/uniqBy'

const SYSTEM_ID = 'bk_cmdb'
const SYSTEM_NAME = '配置平台'

// 前端构造的auth结构为：
// [{ type: 'xxx', relation: [xxx] }]
// 为了便于view中书写，其中relation可能存在三种格式:
// relation: [1, 2, ...] 表示该动作只关联一个视图，relation成员为视图拓扑路径上的资源ID，即关联单视图，操作单资源
// relation: [[1, 2], [3, 4], ...] 表示该动作只关联一个视图，relation中的成员为数组，每个数组代表一个视图的拓扑路径上的资源ID，即关联单视图，操作多资源
// relation: [[[1, 2], [3, 4]], [[1, 2], [5, 6]]] 表示该动作关联两个及以上的视图，为第二种情况的多视图场景，即关联多视图，操作多资源
// 因第一、第二种均为第三种的子场景，因此通过简单的类型判断转换为第三种形式
// 类型判断减少复杂度，只判断第一个元素的类型，不合法的混搭写法会报错
function convertRelation(relation = [], type) {
  if (!relation.length) return relation
  try {
    const [levelOne] = relation
    if (!Array.isArray(levelOne)) { // [1, 2, ...]的场景
      return [[relation]]
    }
    const [levelTwo] = levelOne
    if (!Array.isArray(levelTwo)) {
      return relation.map(data => [data])
    }
    return relation
  } catch (error) {
    $error('Convert resource relations fail, wrong params')
    console.error('Convert resource relations fail, wrong params:')
    console.error('auth type:', type)
    console.error('auth relation:', relation)
  }
}

// 将相同动作下的相同视图的实例合并到一起，并且将相同的实例去重
export function mergeSameActions(actions) {
  const actionMap = new Map()

  actions.forEach((action) => {
    const viewMap = actionMap.get(action.id) || new Map()
    action.related_resource_types.forEach(({ type, instances = [] }) => {
      const viewInstances = viewMap.get(type) || []
      viewInstances.push(...instances)
      viewMap.set(type, viewInstances)
    })
    actionMap.set(action.id, viewMap)
  })

  const permission = {
    system_id: SYSTEM_ID,
    system_name: SYSTEM_NAME,
    actions: []
  }

  actionMap.forEach((viewMap, actionId) => {
    const relatedResourceTypes = []
    viewMap.forEach((viewInstances, viewType) => {
      // 将每个view下的实例去重，viewInstances中每一条实例的结构可能是 [inst] 或者 [inst, inst, ...]，所以必须要合并所有实例以确定其唯一性
      const instances = uniqBy(viewInstances, insts => insts?.reduce((acc, cur) => `${acc}/${cur?.id}_${cur?.type}`, ''))
      const data = {
        type: viewType,
        system_id: SYSTEM_ID,
        system_name: SYSTEM_NAME,
      }
      if (instances?.length) {
        data.instances = instances
      }
      relatedResourceTypes.push(data)
    })
    permission.actions.push({
      id: actionId,
      related_resource_types: relatedResourceTypes
    })
  })
  return permission
}

// 用于转换为无权限申请弹窗中展示需要的数据，结构与接口无权限返回的数据一致，申请权限的接口也是用这个数据
export const translateAuth = (auth) => {
  const authList = Array.isArray(auth) ? auth : [auth]
  const actions = authList.map(({ type, relation = [] }) => {
    relation = convertRelation(relation, type)
    const definition = IAM_ACTIONS[type]
    const action = {
      id: typeof definition.id === 'function' ? definition.id(relation) : definition.id,
      related_resource_types: []
    }
    if (!definition.relation) {
      return action
    }

    // 计算出完整的关联路径用于申请权限展示，如：模型-实例
    definition.relation.forEach((viewDefinition, viewDefinitionIndex) => { // 第m个视图的定义n
      const { view, instances } = viewDefinition
      const relatedResource = {
        type: typeof view === 'function' ? view(relation) : view
      }

      if (instances?.length) {
        relatedResource.instances = []
        relation.forEach((resourceViewPaths) => { // 第x个资源对应的视图数组
          const viewPathData = resourceViewPaths[viewDefinitionIndex] || [] // 取出第x个资源对应的第m个视图对应的拓扑路径ID数组
          if (typeof instances === 'function') {
            relatedResource.instances.push(instances(relation))
          } else {
            const viewFullPath = viewPathData.map((path, pathIndex) => ({ // 资源x的第m个视图对应的全路径拓扑对象
              type: instances[pathIndex],
              id: String(path)
            }))
            if (!relatedResource.instances.some(path => isEqual(path, viewFullPath))) {
              relatedResource.instances.push(viewFullPath)
            }
          }
        })
      }
      action.related_resource_types.push(relatedResource)
    })

    return action
  })
  return mergeSameActions(actions)
}

export function filterPassedAuth(auth, authResults) {
  const authList = Array.isArray(auth) ? auth : [auth]
  return authList.filter(({ type, relation = [] }, index) => {
    // 通用模型实例编辑
    if (OPERATION.U_INST === type) {
      return authResults.filter(item => item.resource_type.startsWith('comobj_')).some((item) => {
        const modelId = item.resource_type?.split('_')?.[1]
        return item.resource_id === relation?.[1] && Number(modelId) === relation?.[0] && !item?.is_pass
      })
    }

    // 业务主机或资源池主机编辑
    if (OPERATION.U_HOST === type || OPERATION.U_RESOURCE_HOST === type) {
      return authResults.some(item => item.resource_id === relation?.[1] && !item?.is_pass)
    }

    if ([OPERATION.U_MODEL, OPERATION.R_MODEL].includes(type)) {
      return !authResults?.[index]?.is_pass
    }

    if (OPERATION.C_FIELD_TEMPLATE === type) {
      return authResults.some(item => item.action === 'create' && !item?.is_pass)
    }

    return !authResults?.[index]?.is_pass
  })
}

cursor.setOptions({
  globalCallback: (options) => {
    const { authResults, ignorePassedAuth, callbackUrl, relatedPermission, showPermissionDialog = true } = options

    // 根据配置去除有权限的auth，在此处去除比在permission中去除会更简单
    let newAuth = options.auth
    if (ignorePassedAuth) {
      newAuth = filterPassedAuth(options.auth, authResults)
    }
    const permission = translateAuth(newAuth)

    const { permissionModal } = window
    showPermissionDialog && permissionModal?.show(permission, authResults, callbackUrl, relatedPermission)
  },
  x: 16,
  y: 8
})
