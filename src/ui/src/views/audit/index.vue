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
  <div class="audit-layout">
    <bk-tab class="audit-tab"
      type="unborder-card"
      :active.sync="active">
      <bk-tab-panel
        v-for="panel in tabPanels"
        v-bind="panel"
        :key="panel.id">
      </bk-tab-panel>
    </bk-tab>
    <div class="audit-options">
      <component ref="optionsComponent" :is="optionsComponent" @condition-change="handleConditionChange"></component>
    </div>
    <bk-table v-bkloading="{ isLoading: $loading(request.list) }"
      :data="table.list"
      :pagination="table.pagination"
      :max-height="$APP.height - 310"
      :row-style="{ cursor: 'pointer' }"
      @page-change="handlePageChange"
      @page-limit-change="handleSizeChange"
      @sort-change="handleSortChange"
      @row-click="handleRowClick">
      <bk-table-column
        prop="resource_type"
        class-name="is-highlight"
        :label="$t('操作对象')"
        :formatter="getResourceTypeName">
      </bk-table-column>
      <bk-table-column
        prop="action"
        :label="$t('动作')"
        :formatter="getActionName">
      </bk-table-column>
      <bk-table-column
        v-if="['host', 'business'].includes(active)"
        prop="bk_biz_name"
        :label="$t('所属业务')">
        <template slot-scope="{ row }">
          <audit-business-selector type="info" :value="row.bk_biz_id"></audit-business-selector>
        </template>
      </bk-table-column>
      <bk-table-column
        prop="resource_name"
        :show-overflow-tooltip="true"
        :label="$t('实例')">
        <template slot-scope="{ row }">{{getResourceName(row)}}</template>
      </bk-table-column>
      <bk-table-column
        :label="$t('操作描述')">
        <template slot-scope="{ row }">{{`${getActionName(row)}${getResourceTypeName(row)}`}}</template>
      </bk-table-column>
      <bk-table-column
        width="160"
        sortable="custom"
        prop="operation_time"
        :label="$t('时间')">
        <template slot-scope="{ row }">{{$tools.formatTime(row.operation_time)}}</template>
      </bk-table-column>
      <bk-table-column
        prop="user"
        :label="$t('操作账号')">
      </bk-table-column>
      <cmdb-table-empty slot="empty" :stuff="table.stuff" @clear="handleClearFilter"></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<script>
  import AuditBusinessOptions from './children/audit-business-options'
  import AuditResourceOptions from './children/audit-resource-options'
  import AuditOtherOptions from './children/audit-other-options'
  import AuditHostOptions from './children/audit-host-options'
  import AuditBusinessSelector from '@/components/audit-history/audit-business-selector'
  import RouterQuery from '@/router/query'
  import AuditDetails from '@/components/audit-history/details.js'
  import { isNumeric } from '@/utils/util'
  export default {
    components: {
      [AuditBusinessOptions.name]: AuditBusinessOptions,
      [AuditResourceOptions.name]: AuditResourceOptions,
      [AuditOtherOptions.name]: AuditOtherOptions,
      [AuditHostOptions.name]: AuditHostOptions,
      AuditBusinessSelector
    },
    data() {
      return {
        active: RouterQuery.get('tab', 'host'),
        tabPanels: [
          {
            name: 'host',
            label: this.$t('主机')
          },
          {
            name: 'business',
            label: this.$t('业务')
          }, {
            name: 'resource',
            label: this.$t('资源')
          }, {
            name: 'other',
            label: this.$t('其他')
          }
        ],
        dictionary: [],
        condition: {},
        table: {
          list: [],
          pagination: this.$tools.getDefaultPaginationConfig({ 'limit-list': [20, 50, 100, 200] }),
          sort: '-operation_time',
          stuff: {
            // 如果初始化时有查询参数，则认为处在查询模式
            type: Object.keys(this.$route.query).length ? 'search' : 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        request: {
          list: Symbol('list')
        }
      }
    },
    computed: {
      optionsComponent() {
        const componentMap = {
          business: AuditBusinessOptions.name,
          resource: AuditResourceOptions.name,
          other: AuditOtherOptions.name,
          host: AuditHostOptions.name
        }
        return componentMap[this.active]
      }
    },
    created() {
      this.setQueryWatcher()
      this.getAuditDictionary()
    },
    beforeDestroy() {
      this.teardownQueryWatcher()
    },
    methods: {
      setQueryWatcher() {
        this.unwatch = RouterQuery.watch('*', ({
          page,
          limit,
          sort,
          tab,
          _e: isEvent
        }) => {
          this.active = tab || 'host'
          this.table.pagination.current = parseInt(page || this.table.pagination.current, 10)
          this.table.pagination.limit = parseInt(limit || this.table.pagination.limit, 10)
          this.table.sort = sort || this.table.sort
          this.$nextTick(() => this.getAuditList(isEvent))
        })
      },
      teardownQueryWatcher() {
        this.unwatch && this.unwatch()
      },
      async getAuditDictionary() {
        try {
          this.dictionary = await this.$store.dispatch('audit/getDictionary', {
            fromCache: true,
            globalPermission: false
          })
        } catch (error) {
          this.dictionary = []
        }
      },
      parseResourceId(id) {
        if (isNumeric(id)) return parseInt(id, 10)
        return id
      },
      handleConditionChange(condition) {
        const useCondition = (condition) => {
          const usefulCondition = {}
          Object.keys(condition).forEach((key) => {
            const value = condition[key]
            if (String(value).length) {
              usefulCondition[key] = value
            }
          })
          // 动态分组的ID是String, 其他的是Number, 区别转换
          if (usefulCondition.resource_id) {
            usefulCondition.resource_id = usefulCondition.resource_type === 'dynamic_group'
              ? usefulCondition.resource_id
              : this.parseResourceId(usefulCondition.resource_id)
          }
          // 转换时间范围为start/end的形式
          if (usefulCondition.operation_time) {
            const [start, end] = usefulCondition.operation_time
            usefulCondition.operation_time = {
              start: `${start} 00:00:00`,
              end: `${end} 23:59:59`
            }
          }
          // 账号多选转为Array
          if (usefulCondition.user) {
            usefulCondition.user = Array.isArray(usefulCondition.user) ? usefulCondition.user : usefulCondition.user.split(',')
          }

          return usefulCondition
        }

        const usefulCondition = useCondition(condition)

        // 处理内联condition
        const { condition: subCondition = {} } = condition
        const usefulSubCondition = []
        Object.keys(subCondition).forEach((key) => {
          const [operator, value] = subCondition[key]
          const conditionValue = useCondition({ [key]: value })
          if (Object.keys(conditionValue).length) {
            usefulSubCondition.push({
              field: key,
              operator,
              value: conditionValue[key]
            })
          }
        })
        usefulCondition.condition = usefulSubCondition

        // 兼容resource_name为in/contains操作
        if (usefulCondition.resource_name || usefulCondition.extend_resource_name) {
          const {
            fuzzy_query: fuzzy,
            resource_name: resourceName,
            extend_resource_name: extendResourceName
          } = usefulCondition
          const field = resourceName ? 'resource_name' : 'extend_resource_name'
          const finalResourceName = resourceName || extendResourceName

          usefulSubCondition.push({
            field,
            operator: fuzzy ? 'contains' : 'in',
            value: fuzzy ? finalResourceName : [finalResourceName]
          })
          delete usefulCondition[field]
          delete usefulCondition.fuzzy_query
        }
        this.condition = usefulCondition
      },
      async getAuditList(eventTrigger) {
        try {
          const params = {
            condition: this.condition,
            page: {
              ...this.$tools.getPageParams(this.table.pagination),
              sort: this.table.sort
            }
          }
          const { info, count } = await this.$store.dispatch('audit/getList', {
            params,
            config: {
              requestId: this.request.list,
              globalPermission: false
            }
          })

          this.table.stuff.type = eventTrigger ? 'search' : 'default'
          this.table.pagination.count = count
          this.table.list = info
        } catch ({ permission }) {
          if (permission) {
            this.$route.meta.view = 'permission'
          }
        }
      },
      handlePageChange(current) {
        RouterQuery.set({
          page: current,
          _t: Date.now()
        })
      },
      handleSizeChange(size) {
        RouterQuery.set({
          limit: size,
          page: 1,
          _t: Date.now()
        })
      },
      handleSortChange(sort) {
        RouterQuery.set({
          page: 1,
          sort: this.$tools.getSort(sort, 'operation_time'),
          _t: Date.now()
        })
      },
      handleRowClick(row) {
        AuditDetails.show({
          id: row.id
        })
      },
      getResourceName(row) {
        return row.resource_name || row.extend_resource_name || '--'
      },
      getResourceTypeName(row) {
        const type = this.dictionary.find(type => type.id === row.resource_type)
        return type ? type.name : row.resource_type
      },
      getActionName(row) {
        const type = this.dictionary.find(type => type.id === row.resource_type)
        const operations = type ? type.operations : []
        const operation = operations.find(operation => operation.id === row.action)
        return operation ? operation.name : row.action
      },
      handleClearFilter() {
        this.$refs.optionsComponent.handleReset()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .audit-layout{
        padding: 0 20px;
        .audit-tab {
            height: auto;
            /deep/ .bk-tab-header {
                padding: 0;
            }
        }
        .filter {
            display: flex;
            align-items: center;
            justify-content: flex-start;
            flex-direction: row;
            flex-wrap: wrap;
            padding: 22px 0 10px 0;
            .option {
                flex: none;
                width: 27%;
                margin: 0 1.5% 12px 0;
                white-space: nowrap;
                font-size: 0;
                display: flex;
                align-items: center;
                &.instance {
                    .bk-select {
                        width: 40%;
                        margin-right: 5px;
                    }
                    .bk-form-control {
                        width: calc(60% - 5px);
                    }
                }

                &.action,
                &.operator {
                    .name {
                        width: 96px;
                        text-align: right;
                    }
                }

                &.resource,
                &.instance {
                    .name {
                        width: 70px;
                        text-align: right;
                    }
                }
            }
            .option-btn {
                width: auto;
                .bk-button + .bk-button {
                    margin-left: 8px;
                }
            }
            .name {
                font-size: 14px;
                padding-right: 10px;
                @include ellipsis;
            }
            .content {
                flex: 1;
                width: calc(100% - 96px);
                .bk-select {
                    width: 100%;
                }
            }
        }

    }

    [bk-language="en"] {
        .filter {
            .name {
                min-width: 70px;
                text-align: right;
            }

            .option {
                &.action,
                &.operator {
                    .name {
                        width: 146px;
                        text-align: right;
                    }
                }

                &.resource,
                &.instance {
                    .name {
                        width: 130px;
                        text-align: right;
                    }
                }
            }
        }
    }
</style>
