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

import $http from '@/api'
import store from '@/store'
import { BUILTIN_MODELS } from '@/dictionary/model-constants'

/**
 * @typedef ParentLayers 父级依赖关系
 * @property {string} resource_type 父级依赖，例如实例依赖模型
 * @property {string} resource_id 父级依赖id
 */

/**
 * @typedef Resource 权限信息
 * @property {string} action 需要鉴权的操作
 * @property {string} resource_type 需要鉴权的资源类型
 * @property {string} [bk_biz_id] 对业务下的操作进行鉴权需要提供业务 ID
 * @property {string} [resource_id] 资源 ID
 * @property {Array.<ParentLayers>} [parent_layers]
 */

/**
 * 鉴定用户是否有某资源的某操作的权限
 * @param {Array.<Resource>} resources 权限信息列表
 * @returns {Promise}
 */
export const verifyAuth = resources => $http.post('auth/verify', {
  resources
})

/**
 * 获取免查看鉴权的模型
 */
export const getViewAuthFreeModels = () => store.getters['objectModelClassify/viewAuthFreeModels']

/**
 * 判断一个模型是否为免查看鉴权
 * @param {Object} model 单个模型
 */
export const isViewAuthFreeModel = (model) => {
  const dataKey = model.bk_obj_id ? 'bk_obj_id' : 'id'
  return store.getters['objectModelClassify/viewAuthFreeModelSet'].has(model[dataKey])
}

/**
 * 根据模型判断其实例是否为免查看鉴权
 * @param {Object} model 单个模型
 */
export const isViewAuthFreeModelInstance = (model) => {
  const models = store.getters['objectModelClassify/models']
  const authFreeModelInstances = [
    BUILTIN_MODELS.BUSINESS,
    BUILTIN_MODELS.BUSINESS_SET,
    BUILTIN_MODELS.HOST,
    BUILTIN_MODELS.SET,
    BUILTIN_MODELS.MODULE,
    BUILTIN_MODELS.PROJECT
  ]
  const dataKey = model.bk_obj_id ? 'bk_obj_id' : 'id'
  let objId = model[dataKey]
  if (dataKey === 'id') {
    objId = models.find(item => item.id === model.id)?.bk_obj_id
  }
  return authFreeModelInstances.includes(objId)
}
