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

import { reactive, toRefs, set, del } from 'vue'
import { t } from '@/i18n'
import useClone from '@/hooks/utils/clone'
const defaultState = {
  visible: false,
  title: '',
  bk_obj_id: null,
  bk_biz_id: null,
  available: () => true,
  submit: () => {},
  count: 0,
  limit: 30000,
  step: 1,
  steps: [{ title: t('选择字段'), icon: 1 }, { title: t('选择关联模型'), icon: 2 }],
  status: null,
  presetFields: [],
  defaultSelectedFields: [],
  fields: [],
  originFields: [],
  relations: {},
  exportRelation: false,
  object_unique_id: '',
  confirmBtnText: t('开始导出')
}

const state = reactive(useClone(defaultState))

const setState = (newState) => {
  Object.assign(state, newState)
}

const resetState = () => setState(useClone(defaultState))
const resetPartial = () => setState({
  step: 1,
  status: null,
  relations: {},
  exportRelation: false,
  object_unique_id: ''
})

const setRelation = (modelId, uniqueId) => set(state.relations, modelId, uniqueId)
const removeRelation = modelId => del(state.relations, modelId)

export default function () {
  return [toRefs(state), { setState, resetState, resetPartial, setRelation, removeRelation }]
}
