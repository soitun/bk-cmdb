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
const actions = {
  getList(context, { params, config }) {
    return $http.post('findmany/audit_list', params, config)
  },
  getDictionary(context, config) {
    return $http.get('find/audit_dict', config)
  },
  getDetails(context, { id, config }) {
    return $http.post('find/audit', { id: [id] }, config).then(([detail]) => detail)
  },
  getInstList(context, { params, config }) {
    return $http.post('find/inst_audit', params, config)
  },
  getInstDetails(context, { params, config }) {
    return $http.post('find/inst_audit', params, config).then(({ info: [detail] }) => detail)
  }
}

export default {
  namespaced: true,
  actions
}
