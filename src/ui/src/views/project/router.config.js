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
import { MENU_RESOURCE_PROJECT, MENU_RESOURCE_PROJECT_DETAILS, MENU_RESOURCE_MANAGEMENT } from '@/dictionary/menu-symbol'
import { OPERATION } from '@/dictionary/iam-auth'

export default [
  {
    name: MENU_RESOURCE_PROJECT,
    path: 'project',
    component: () => import('./index.vue'),
    meta: new Meta({
      menu: {
        i18n: '项目',
        relative: MENU_RESOURCE_MANAGEMENT
      },
      auth: {
        view: { type: OPERATION.R_PROJECT }
      }
    })
  },
  {
    name: MENU_RESOURCE_PROJECT_DETAILS,
    path: 'project/details/:projId',
    component: () => import('./details.vue'),
    meta: new Meta({
      menu: {
        i18n: '项目详情',
        relative: MENU_RESOURCE_PROJECT
      },
      layout: {},
      auth: {
        view: { type: OPERATION.R_PROJECT }
      }
    })
  }
]
