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
  <bk-date-picker
    type="daterange"
    transfer
    :value="localValue"
    v-bind="$attrs"
    format="yyyy-MM-dd"
    @change="handleChange"
    @open-change="handleToggle"
    @clear="() => $emit('clear')">
  </bk-date-picker>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-date',
    mixins: [activeMixin],
    props: {
      value: {
        type: Array,
        default: () => ([])
      }
    },
    computed: {
      localValue: {
        get() {
          return [...this.value]
        },
        set(values) {
          this.$emit('input', values)
          this.$emit('change', values)
        }
      }
    },
    methods: {
      handleChange(values) {
        if (values.toString() === this.value.toString()) return
        this.localValue = values.filter(value => !!value)
      }
    }
  }
</script>
