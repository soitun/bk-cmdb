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
  <div class="process-bind-info-value">
    <bk-popover v-bind="popoverOptoins" :disabled="popoverList.length < 2">
      <table-value
        ref="table"
        :value="localValue"
        :show-on="'cell'"
        :format-cell-value="formatCellValue"
        :property="property">
      </table-value>
      <ul slot="content">
        <li v-for="(item, index) in popoverList" :key="index">{{item}}</li>
      </ul>
    </bk-popover>
  </div>
</template>

<script>
  import TableValue from '@/components/ui/other/table-value'
  import { PROCESS_BIND_IP_ALL_MAP } from '@/dictionary/process-bind-ip.js'

  export default {
    components: {
      TableValue
    },
    props: {
      value: {
        type: Array,
        default: () => ([])
      },
      property: {
        type: Object,
        default: () => ({})
      },
      popoverOptoins: {
        type: Object,
        default: () => ({})
      }
    },
    data() {
      return {
        popoverList: [],
        localValue: []
      }
    },
    watch: {
      value: {
        handler(value) {
          this.localValue = value || []
          this.setPopoverList()
        },
        immediate: true
      }
    },
    methods: {
      ipText(value) {
        const map = PROCESS_BIND_IP_ALL_MAP
        return map[value] || value || '--'
      },
      setPopoverList() {
        this.$nextTick(() => {
          const list = this.$refs.table.cellValue
          this.popoverList = list.map(this.getRowValue)
        })
      },
      getRowValue(row) {
        const ip = this.ipText(row.ip)
        return `${row.protocol} ${ip}:${row.port}`
      },
      formatCellValue(list) {
        if (!list.length) {
          return '--'
        }
        const newList = list.map(this.getRowValue)
        const total = list.length
        const showCount = total > 1
        return (
          <div class={`bind-info-value${showCount ? ' show-count' : ''}`}>
            <span>{newList.join(', ')}</span>
            {showCount ? <span class="count">{total}</span> : ''}
          </div>
        )
      }
    }
  }
</script>

<style lang="scss" scoped>
    .process-bind-info-value {
        /deep/ .bind-info-value {
            position: relative;
            display: inline-block;
            vertical-align: middle;
            padding-right: 0;
            max-width: 100%;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;

            &.show-count {
                padding-right: 26px;
            }

            .count {
                position: absolute;
                display: inline-block;
                right: 2px;
                top: 0;
                color: #979ba5;
                background-color: #f0f1f5;
                border-radius: 3px;
                padding: 0 2px;
            }
        }
        /deep/ .bk-tooltip {
            display: block;
        }
        /deep/ .bk-tooltip-ref {
            display: block;
        }
    }
</style>
