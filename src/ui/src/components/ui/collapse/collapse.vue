<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div :class="['collapse-layout', size]">
    <div class="collapse-trigger" @click="toggle">
      <span :class="['collapse-arrow', 'bk-icon', arrowIconClass, { 'is-collapsed': hidden }]"></span>
      <span class="collapse-text" v-bk-overflow-tips>
        <slot name="title">{{label}}</slot>
      </span>
    </div>
    <cmdb-collapse-transition
      @before-enter="handleBeforeEnter"
      @enter="handleEnter"
      @after-enter="handleAfterEnter"
      @enter-cancelled="handleEnterCancelled"
      @before-leave="handleBeforeLeave"
      @leave="handleLeave"
      @after-leave="handleAfterLeave"
      @leave-cancelled="handleLeaveCancelled">
      <div class="collapse-content" v-show="!hidden">
        <slot></slot>
      </div>
    </cmdb-collapse-transition>
  </div>
</template>

<script>
  import isEqual from 'lodash/isEqual'
  export default {
    name: 'cmdb-collapse',
    props: {
      collapse: Boolean,
      label: {
        type: String
      },
      arrowType: {
        type: String,
        default: 'outlined', // filled
      },
      size: {
        type: String
      },
      autoExpand: { // 是否可以自动展开
        type: Boolean,
        default: false
      },
      list: {
        type: [Object, Array],
        default: (() => {})
      }
    },
    data() {
      return {
        hidden: this.collapse
      }
    },
    computed: {
      arrowIconClass() {
        const classMap = {
          outlined: 'icon-angle-down',
          filled: 'icon-down-shape'
        }
        return `${classMap[this.arrowType]} ${this.arrowType}`
      }
    },
    watch: {
      list(val, lastVal) {
        const isSame = isEqual(val, lastVal)
        if (!isSame && this.hidden && this.autoExpand) {
          this.toggle()
        }
      },
      collapse(collapse) {
        this.hidden = collapse
      },
      hidden(hidden) {
        this.$emit('update:collapse', hidden)
        this.$emit('collapse-change', hidden)
      }
    },
    methods: {
      toggle() {
        this.hidden = !this.hidden
      },
      handleBeforeEnter() {
        this.$emit('before-enter')
      },
      handleEnter() {
        this.$emit('enter')
      },
      handleAfterEnter() {
        this.$emit('after-enter')
      },
      handleEnterCancelled() {
        this.$emit('enter-cancelled')
      },
      handleBeforeLeave() {
        this.$emit('before-leave')
      },
      handleLeave() {
        this.$emit('leave')
      },
      handleAfterLeave() {
        this.$emit('after-leave')
      },
      handleLeaveCancelled() {
        this.$emit('leave-cancelled')
      }
    }
  }
</script>

<style lang="scss">
    .collapse-layout {
        .collapse-trigger {
            display: flex;
            color: #333948;
            font-weight: bold;
            align-items: center;
            cursor: pointer;
            .collapse-arrow {
                font-size: 20px;
                font-weight: 700;
                margin: 0 2px 0 -4px;
                transition: transform .2s ease-in-out;
                &.is-collapsed {
                    transform: rotate(-90deg);
                }

                &.filled {
                  font-size: 12px;
                  color: #63656E;
                  margin: 0 4px 0 0;
                }
            }
            .collapse-text {
                flex: 1;
                font-size: 14px;
                @include ellipsis;
            }
        }

        &.small {
            .collapse-arrow {
                &.filled {
                  margin-top: -1px;
                }
            }
            .collapse-text {
                font-size: 12px;
            }
        }
    }
</style>
