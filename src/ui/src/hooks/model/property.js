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

import { reactive, toRefs, watch, isRef } from 'vue'
import { BUILTIN_MODELS } from '@/dictionary/model-constants'
import propertyService from '@/service/property/property'
export default function (options = {}, config) {
  const state = reactive({
    properties: [],
    invisibleProperties: [],
    pending: false
  })
  const refresh = async (value) => {
    if (!value.bk_obj_id) return
    state.pending = true
    if (value.bk_obj_id === BUILTIN_MODELS.PROJECT) {
      state.invisibleProperties = ['bk_project_icon']
    }
    state.properties = (await propertyService.find(value, config))
      .filter(property => !state.invisibleProperties.includes(property.bk_property_id))
    state.pending = false
  }
  watch(() => (isRef(options) ? options.value : options), refresh, { immediate: true, deep: true })
  return [toRefs(state), { refresh }]
}
