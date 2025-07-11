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

import Meta from '@/router/meta'
import { MENU_MODEL_MANAGEMENT, MENU_MODEL_DETAILS } from '@/dictionary/menu-symbol'
import { OPERATION } from '@/dictionary/iam-auth'

export default [
  {
    name: MENU_MODEL_MANAGEMENT,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
      menu: {
        i18n: '模型管理'
      }
    })
  },
  {
    name: MENU_MODEL_DETAILS,
    path: 'index/details/:modelId',
    component: () => import('./children/model-details/index.vue'),
    meta: new Meta({
      menu: {
        i18n: '模型详情',
        relative: MENU_MODEL_MANAGEMENT
      },
      layout: {},
      checkAvailable: (to, from, app) => {
        const { modelId } = to.params
        const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
        return !!model
      },
      auth: {
        view: (to, app) => {
          const { modelId } = to.params
          const model = app.$store.getters['objectModelClassify/getModelById'](modelId)
          return ({ type: OPERATION.R_MODEL, relation: [model.id] })
        }
      }
    })
  }
]
