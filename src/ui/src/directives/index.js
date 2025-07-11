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

import Vue from 'vue'
import vClickOutside from 'v-click-outside'
import cursor from './cursor.js'
import transferDom from './transfer-dom.js'
import setTestId from './set-test-id.js'
import { autofocus } from './autofocus.js'
import { scroll } from './scroll.js'

Vue.use(vClickOutside)
Vue.use(cursor)
Vue.use(setTestId)
Vue.directive('autofocus', autofocus)
Vue.directive('transfer-dom', transferDom)
Vue.directive('scroll', scroll)

export default {
  'v-click-outside': vClickOutside,
  'v-transfer-dom': transferDom
}
