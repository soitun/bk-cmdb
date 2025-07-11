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

/**
 * @directive 自动聚焦指令，目前只支持 input
 */
const focusCall = (el) => {
  const input = el.querySelector('input') || el.querySelector('textarea')
  if (input) {
    // 尽量靠后执行，避免其他队列任务影响到聚焦
    setTimeout(() => {
      input.focus()
      // 将光标放在末尾
      const length = input?.value?.length ?? 0
      input.setSelectionRange(length, length)
    }, 0)
  }
}
export const autofocus = {
  update: (el, { modifiers }) => {
    if (!modifiers?.update) return
    focusCall(el)
  },
  inserted: (el) => {
    focusCall(el)
  }
}
