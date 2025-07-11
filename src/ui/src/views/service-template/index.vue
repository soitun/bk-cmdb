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
  <div class="template-wrapper" ref="templateWrapper">
    <cmdb-tips class="mb10 top-tips" tips-key="serviceTemplateTips">
      <i18n path="服务模板功能提示">
        <template #link>
          <a class="tips-link" href="javascript:void(0)" @click="handleTipsLinkClick">{{$t('业务拓扑')}}</a>
        </template>
      </i18n>
    </cmdb-tips>
    <div class="template-filter clearfix">
      <cmdb-auth class="fl mr10" :auth="{ type: $OPERATION.C_SERVICE_TEMPLATE, relation: [bizId] }">
        <bk-button slot-scope="{ disabled }" v-test-id="'create'"
          theme="primary"
          :disabled="disabled"
          @click="handleCreate">
          {{$t('新建')}}
        </bk-button>
      </cmdb-auth>
      <div class="filter-text fr">
        <bk-select class="fl"
          font-size="medium"
          :placeholder="$t('所有一级分类')"
          :allow-clear="true"
          :searchable="true"
          v-model="filter.mainClassification"
          @selected="handleSelect"
          @clear="() => handleSelect()">
          <bk-option v-for="category in mainList"
            :key="category.id"
            :id="category.id"
            :name="category.name">
          </bk-option>
        </bk-select>
        <bk-select class="fl"
          font-size="medium"
          :placeholder="$t('所有二级分类')"
          :allow-clear="true"
          :searchable="true"
          :empty-text="emptyText"
          v-model="filter.secondaryClassification"
          @selected="handleSelectSecondary"
          @clear="() => handleSelectSecondary()">
          <bk-option v-for="category in secondaryList"
            :key="category.id"
            :id="category.id"
            :name="category.name">
          </bk-option>
        </bk-select>
        <bk-input type="text"
          class="filter-search fl"
          :placeholder="$t('请输入xx', { name: $t('模板名称') })"
          :right-icon="'bk-icon icon-search'"
          clearable
          font-size="medium"
          v-model.trim="filter.templateName"
          @enter="setRoute()"
          @right-icon-click="setRoute()"
          @clear="handlePageChange(1)">
        </bk-input>
      </div>
    </div>
    <bk-table class="template-table" v-test-id="'templateList'"
      v-bkloading="{ isLoading: $loading(request.list) }"
      :data="table.list"
      :pagination="table.pagination"
      :max-height="$APP.height - 229"
      :row-style="{ cursor: 'pointer' }"
      @row-click="handleRowClick"
      @page-limit-change="handleSizeChange"
      @page-change="handlePageChange"
      @sort-change="handleSortChange">
      <bk-table-column prop="id" label="ID" class-name="is-highlight" sortable="custom">
        <div slot-scope="{ row }" v-bk-overflow-tips
          :class="['template-id', { 'need-sync': row.need_sync }]">
          {{row.id}}
        </div>
      </bk-table-column>
      <bk-table-column prop="name" :label="$t('模板名称')" show-overflow-tooltip sortable="custom"></bk-table-column>
      <bk-table-column prop="service_category" :label="$t('服务分类')" show-overflow-tooltip></bk-table-column>
      <bk-table-column prop="process_template_count" :label="$t('进程数量')">
        <template slot-scope="{ row }">
          <cmdb-loading :loading="$loading(request.count)">
            <template v-if="row.process_template_count > 0">
              {{row.process_template_count}}
            </template>
            <span style="color: #ff9c01" v-else>{{row.process_template_count}}（{{$t('未配置')}}）</span>
          </cmdb-loading>
        </template>
      </bk-table-column>
      <bk-table-column prop="module_count" :label="$t('已应用模块数')">
        <template slot-scope="{ row }">
          <cmdb-loading :loading="$loading(request.count)">{{row.module_count}}</cmdb-loading>
        </template>
      </bk-table-column>
      <bk-table-column prop="modifier" :label="$t('修改人')" sortable="custom"></bk-table-column>
      <bk-table-column prop="last_time" :label="$t('修改时间')" show-overflow-tooltip sortable="custom">
        <template slot-scope="{ row }">
          {{$tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm')}}
        </template>
      </bk-table-column>
      <bk-table-column prop="operation" :label="$t('操作')" fixed="right">
        <template slot-scope="{ row }">
          <cmdb-loading :loading="$loading(request.count)">
            <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_SERVICE_TEMPLATE, relation: [bizId, row.id] }">
              <bk-button slot-scope="{ disabled }"
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleEdit(row.id)">
                {{$t('编辑')}}
              </bk-button>
            </cmdb-auth>
            <cmdb-auth class="mr10" :auth="{ type: $OPERATION.C_SERVICE_TEMPLATE, relation: [bizId] }">
              <bk-button slot-scope="{ disabled }"
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="cloneTemplate(row.id)">
                {{$t('克隆')}}
              </bk-button>
            </cmdb-auth>
            <cmdb-auth :auth="{ type: $OPERATION.D_SERVICE_TEMPLATE, relation: [bizId, row.id] }">
              <template slot-scope="{ disabled }">
                <span class="text-primary"
                  style="color: #dcdee5 !important; cursor: not-allowed;"
                  v-if="row['module_count'] && !disabled"
                  v-bk-tooltips.top="$t('不可删除')">
                  {{$t('删除')}}
                </span>
                <bk-button v-else v-test-id="'delTemplate'"
                  theme="primary"
                  :disabled="disabled"
                  :text="true"
                  @click.stop="deleteTemplate(row)">
                  {{$t('删除')}}
                </bk-button>
              </template>
            </cmdb-auth>
          </cmdb-loading>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff"
        :auth="{ type: $OPERATION.C_SERVICE_TEMPLATE, relation: [bizId] }"
        @create="handleCreate"
        @clear="handleClearFilter"
      ></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<script>
  import { mapActions, mapGetters } from 'vuex'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SERVICE_TEMPLATE_CREATE,
    MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
    MENU_BUSINESS_SERVICE_TEMPLATE_EDIT
  } from '@/dictionary/menu-symbol'
  import CmdbLoading from '@/components/loading/loading'
  import RouterQuery from '@/router/query'
  export default {
    components: {
      CmdbLoading
    },
    data() {
      return {
        filter: {
          mainClassification: '',
          secondaryClassification: '',
          templateName: ''
        },
        table: {
          height: 600,
          list: [],
          pagination: {
            current: 1,
            count: 0,
            ...this.$tools.getDefaultPaginationConfig({
              'limit-list': [20, 50, 100]
            })
          },
          sort: '-id',
          stuff: {
            type: 'default',
            payload: {
              resource: this.$t('服务模板')
            }
          }
        },
        mainList: [],
        secondaryList: [],
        allSecondaryList: [],
        originTemplateData: [],
        maincategoryId: null,
        categoryId: null,
        request: {
          list: Symbol('list'),
          count: Symbol('count')
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      query() {
        return RouterQuery.getAll()
      },
      params() {
        const maincategoryId = this.maincategoryId ? this.maincategoryId : 0
        const id = this.categoryId
          ? this.categoryId
          : maincategoryId
        return {
          bk_biz_id: this.bizId,
          service_category_id: id,
          search: this.filter.templateName,
          page: {
            start: (this.table.pagination.current - 1) * this.table.pagination.limit,
            limit: this.table.pagination.limit,
            sort: this.table.sort
          }
        }
      },
      emptyText() {
        return this.filter.mainClassification ? this.$t('没有二级分类') : this.$t('请选择一级分类')
      },
      hasFilter() {
        return Object.values(this.filter).some(value => !!value)
      }
    },
    watch: {
      query() {
        this.getQueryList()
      }
    },
    async created() {
      try {
        await this.getServiceClassification()
        // 初始化回显选择框数据
        this.initSelect()
        this.getQueryList()
      } catch (e) {
        console.log(e)
      }
    },
    methods: {
      ...mapActions('serviceTemplate', [
        'searchServiceTemplate',
        'deleteServiceTemplate',
        'getServiceTemplateSyncStatus'
      ]),
      ...mapActions('serviceClassification', ['searchServiceCategoryWithoutAmout']),
      getQueryList() {
        const { current = 1, limit = 20, sort = '-id', mainClassification = '', secondaryClassification = '', name = '' } = this.query
        this.table.pagination.current = parseInt(current, 10)
        this.table.pagination.limit = parseInt(limit, 10)
        this.table.sort = sort
        this.filter = {
          mainClassification,
          secondaryClassification,
          templateName: name
        }
        this.getTableData()
      },
      setRoute() {
        const { sort, pagination } = this.table
        const { current, limit } = pagination
        const { mainClassification, secondaryClassification, templateName } = this.filter
        RouterQuery.set({ sort, current, limit,
                          mainClassification, secondaryClassification, name: templateName,
                          _t: Date.now() })
      },
      async getTableData() {
        try {
          const templateData = await this.getTemplateData()
          if (templateData.count && !templateData.info.length) {
            this.table.pagination.current -= 1
            this.setRoute()
          }
          this.table.pagination.count = templateData.count
          this.table.list = templateData.info.map((template) => {
            const second = this.allSecondaryList.find(clazz => clazz.id === template.service_category_id)
            const mainCategory = this.mainList.find(clazz => second && clazz.id === second.bk_parent_id)
            const secondName = second ? second.name : '--'
            const mainCategoryName = mainCategory ? mainCategory.name : '--'
            template.service_category = `${mainCategoryName} / ${secondName}`
            return template
          })
          this.table.stuff.type = this.hasFilter ? 'search' : 'default'
          if (this.table.list.length) {
            this.getTemplateCount()
            this.getTemplateSyncStatus()
          }
        } catch ({ permission }) {
          if (permission) {
            this.table.stuff = {
              type: 'permission',
              payload: { permission }
            }
          }
        }
      },
      async getTemplateCount() {
        try {
          const data = await this.$store.dispatch('serviceTemplate/searchServiceTemplateCount', {
            bizId: this.bizId,
            params: {
              service_template_ids: this.table.list.map(row => row.id)
            },
            config: {
              requestId: this.request.count,
              cancelPrevious: true
            }
          })
          this.table.list.forEach((row) => {
            const counts = data.find(counts => counts.service_template_id === row.id) || {}
            const {
              module_count: moduleCount = '--',
              process_template_count: processTemplateCount = '--'
            } = counts
            this.$set(row, 'module_count', moduleCount)
            this.$set(row, 'process_template_count', processTemplateCount)
          })
        } catch (error) {
          console.error(error)
          this.table.list.forEach((row) => {
            this.$set(row, 'module_count', '--')
            this.$set(row, 'process_template_count', '--')
          })
        }
      },
      getTemplateData() {
        return this.searchServiceTemplate({
          params: this.params,
          config: {
            requestId: this.request.list,
            cancelPrevious: true,
            globalPermission: false
          }
        })
      },
      async getServiceClassification() {
        const { info: categories } = await this.searchServiceCategoryWithoutAmout({
          params: { bk_biz_id: this.bizId },
          config: {
            requestId: 'get_proc_services_categories'
          }
        })
        this.classificationList = categories
        this.mainList = this.classificationList.filter(classification => !classification.bk_parent_id)
        this.allSecondaryList = this.classificationList.filter(classification => classification.bk_parent_id)
      },
      async getTemplateSyncStatus() {
        try {
          const { service_templates: syncStatusList = [] } = await this.getServiceTemplateSyncStatus({
            bizId: this.bizId,
            params: {
              is_partial: true,
              service_template_ids: this.table.list.map(row => row.id)
            },
            config: {
              cancelPrevious: true
            }
          })
          this.table.list.forEach((row) => {
            const syncStatus = syncStatusList.find(status => status.service_template_id === row.id) || {}
            if (syncStatus) {
              this.$set(row, 'need_sync', syncStatus.need_sync)
            }
          })
        } catch (error) {
          console.error(error)
          this.table.list.forEach((row) => {
            this.$set(row, 'need_sync', false)
          })
        }
      },
      initSelect() {
        const { mainClassification = '', secondaryClassification = '' } = this.query
        this.setSelectId(parseInt(mainClassification, 10) || '', parseInt(secondaryClassification, 10) || '')
      },
      setSelectId(id = '', secondId = '') {
        this.secondaryList = this.allSecondaryList.filter(classification => classification.bk_parent_id === id)
        this.maincategoryId = id
        if (secondId) {
          this.categoryId = secondId
        }
      },
      handleSelect(id = '') {
        this.setSelectId(id)
        this.handleSelectSecondary()
      },
      handleSelectSecondary(id = '') {
        this.categoryId = id
        this.filter.secondaryClassification = id
        this.setRoute()
      },
      cloneTemplate(sourceTemplateId) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_CREATE,
          query: {
            clone: sourceTemplateId
          },
          history: true
        })
      },
      handleCreate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_CREATE,
          history: true
        })
      },
      handleEdit(templateId) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_EDIT,
          params: {
            templateId
          },
          history: true
        })
      },
      deleteTemplate(template) {
        this.$bkInfo({
          title: this.$t('确认删除模板'),
          extCls: 'bk-dialog-sub-header-center',
          confirmFn: async () => {
            await this.deleteServiceTemplate({
              params: {
                data: {
                  bk_biz_id: this.bizId,
                  service_template_id: template.id
                }
              },
              config: {
                requestId: 'delete_proc_service_template'
              }
            }).then(() => {
              this.$success(this.$t('删除成功'))
              this.getTableData()
            })
          }
        })
      },
      handleSortChange(sort) {
        this.table.sort = this.$tools.getSort(sort, '-id')
        this.handlePageChange(1)
      },
      handleSizeChange(size) {
        this.table.pagination.limit = size
        this.handlePageChange(1)
      },
      handlePageChange(page) {
        this.table.pagination.current = page
        this.setRoute()
      },
      handleRowClick(row, event, column) {
        if (column.property === 'operation') return
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
          params: {
            templateId: row.id
          },
          history: true
        })
      },
      handleTipsLinkClick() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE
        })
      },
      async handleClearFilter() {
        this.filter.mainClassification = ''
        this.filter.secondaryClassification = ''
        this.filter.templateName = ''
        this.maincategoryId = 0
        this.categoryId = null
        this.secondaryList = []
        await this.getServiceClassification()
        this.setRoute()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .template-wrapper {
        padding: 15px 20px 0;
        .tips-link {
            color: #3A84FF;
            margin: 0;
        }
        .filter-text {
            .bk-select {
                width: 184px;
                margin-right: 10px;
            }
            .filter-search {
                width: 210px;
                position: relative;
            }
        }
        .template-table {
            margin-top: 14px;
            .template-id {
                position: relative;
                display: inline-block;
                vertical-align: middle;
                line-height: 40px;
                padding-right: 10px;
                max-width: 100%;
                @include ellipsis;
                &.need-sync:after {
                    content: "";
                    position: absolute;
                    top: 10px;
                    right: 2px;
                    width: 6px;
                    height: 6px;
                    border-radius: 50%;
                    background-color: $dangerColor;
                }
            }
        }
    }
</style>
