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

/* eslint-disable no-unused-vars */

import $http from '@/api'

const state = {

}

const getters = {

}

const actions = {
  /**
     * 导出模型属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 模型id
     * @return {promises} promises 对象
     */
  exportObjectAttribute({ commit, state, dispatch, rootGetters }, { objId, params, config }) {
    return $http.post(`${window.API_HOST}object/owner/${rootGetters.supplierAccount}/object/${objId}/export`, params, config)
  },

  /**
     * 导入模型属性
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {String} bkObjId 模型id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
  importObjectAttribute({ commit, state, dispatch, rootGetters }, { objId, params, config }) {
    return $http.post(`${window.API_HOST}object/owner/${rootGetters.supplierAccount}/object/${objId}/import`, params, config)
  }
}

const mutations = {

}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
