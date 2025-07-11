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
  <div :class="['template-tree', mode]" v-bkloading="{ isLoading: $loading(request.topo) }">
    <div class="node-root clearfix">
      <i class="folder-icon bk-icon icon-down-shape fl"></i>
      <i class="node-icon fl">{{setName[0]}}</i>
      <span class="root-name" :title="templateName">{{templateName}}</span>
    </div>
    <ul class="node-children">
      <li class="node-child clearfix"
        v-for="(service, index) in services"
        :key="service.id"
        :class="{ selected: selected === service.id }"
        @click="handleChildClick(service)">
        <i class="node-icon fl">{{moduleName[0]}}</i>
        <span class="child-options fr" v-if="mode !== 'view'">
          <bk-link class="action-link" @click="handleViewService(service)" theme="primary">{{$t('查看详情')}}</bk-link>
          <bk-popover v-if="serviceExistHost(service.id)">
            <bk-link class="action-link disabled" theme="primary">{{$t('删除')}}</bk-link>
            <i18n path="该模块下有主机不可删除" tag="p" class="service-tips" slot="content">
              <template #link><span @click="handleGoTopoBusiness(service)">{{$t('跳转查看')}}</span></template>
            </i18n>
          </bk-popover>
          <bk-link v-else class="action-link" @click="handleDeleteService(index)" theme="primary">{{$t('删除')}}</bk-link>
        </span>
        <span class="child-name" v-bk-overflow-tips>{{service.name}}</span>
      </li>
      <li class="options-child node-child clearfix"
        v-if="['create', 'edit'].includes(mode)"
        @click="handleAddService">
        <i class="node-icon icon icon-cc-zoom-in fl"></i>
        <span class="child-name">{{$t('添加服务模板')}}</span>
      </li>
    </ul>
    <bk-dialog
      header-position="left"
      :draggable="false"
      :mask-close="dialog.name !== 'add'"
      :width="840"
      :title="dialog.title"
      v-model="dialog.visible"
      @after-leave="handleDialogClose">
      <component
        ref="dialogComponent"
        :is="dialog.component"
        :services-host="servicesHost"
        v-bind="dialog.props"
        @select-change="handleSelectChange"
        @template-loaded="handleServiceTemplateLoaded">
      </component>
      <template slot="footer">
        <div class="dialog-footer" v-if="dialog.name === 'add'">
          <div class="summary" v-if="serviceTemplateCount > 0">
            <i18n path="已选个数">
              <template #count>
                <span class="stat">
                  <em class="num">{{selectedServiceCount}}</em>
                </span></template>
            </i18n>
            <bk-link class="to-template" theme="primary" icon="icon-cc-share" @click="handleLinkClick">
              {{$t('跳转服务模板')}}
            </bk-link>
          </div>
          <div class="action">
            <bk-button theme="primary" class="btn" @click.stop="handleDialogConfirm">{{$t('确定')}}</bk-button>
            <bk-button theme="default" class="btn ml5" @click.stop="dialog.visible = false">{{$t('取消')}}</bk-button>
          </div>
        </div>
        <bk-button v-else @click="dialog.visible = false">{{$t('关闭')}}</bk-button>
      </template>
    </bk-dialog>
  </div>
</template>

<script>
  import { MENU_BUSINESS_HOST_AND_SERVICE, MENU_BUSINESS_SERVICE_TEMPLATE } from '@/dictionary/menu-symbol'
  import serviceTemplateSelector from './service-template-selector.vue'
  import serviceTemplateInfo from './service-template-info.vue'
  import { rollReqByDataKey } from '@/service/utils'

  export default {
    components: {
      serviceTemplateSelector,
      serviceTemplateInfo
    },
    props: {
      mode: {
        type: String,
        required: true
      },
      templateId: {
        type: Number
      }
    },
    data() {
      return {
        templateName: this.$t('模板集群名称'),
        services: [],
        originalServices: [],
        selected: null,
        unwatch: null,
        dialog: {
          visible: false,
          title: '',
          props: {},
          name: ''
        },
        servicesHost: [],
        selectedServiceCount: 0,
        serviceTemplateCount: 0,
        request: {
          topo: Symbol('topo')
        }
      }
    },
    computed: {
      hasChange() {
        if (this.mode !== 'edit') {
          return false
        }
        if (this.originalServices.length !== this.services.length) {
          return true
        }
        return this.originalServices.some((service, index) => {
          const target = this.services[index]
          return (target && target.id !== service.id) || !target
        })
      },
      setName() {
        const setModel = this.$store.getters['objectModelClassify/getModelById']('set') || {}
        return setModel.bk_obj_name || ''
      },
      moduleName() {
        const moduleModel = this.$store.getters['objectModelClassify/getModelById']('module') || {}
        return moduleModel.bk_obj_name || ''
      },
      sortedServices() {
        return [...this.services].sort((A, B) => A.name.localeCompare(B.name, 'zh-Hans-CN', { sensitivity: 'accent' }))
      },
      topoNodes() {
        return this.originalServices.map(service => ({
          bk_obj_id: 'module',
          bk_inst_id: service.id
        }))
      }
    },
    watch: {
      hasChange(value) {
        this.$emit('service-change', value)
      },
      services(value) {
        this.$emit('service-selected', value)
      },
      mode() {
        this.selected = null
      }
    },
    created() {
      this.initMonitorTemplateName()
      if (['edit', 'view'].includes(this.mode)) {
        this.getSetTemplateServices()
      }
    },
    beforeDestory() {
      this.unwatch && this.unwatch()
    },
    methods: {
      async getSetTemplateServices() {
        try {
          const data = await this.$store.dispatch('setTemplate/getSetTemplateServices', {
            bizId: this.$store.getters['objectBiz/bizId'],
            setTemplateId: this.templateId,
            config: { requestId: this.request.topo }
          })
          this.services = data || []
          this.originalServices = [...this.services]

          if (this.services?.length) {
            this.getServicesHost()
          }
        } catch (e) {
          console.error(e)
          this.services = []
          this.originalServices = []
          this.servicesHost = []
        }
      },
      async getServicesHost() {
        const serviceTemplateIds = this.services.map(item => item.id)
        const hostCounts = await rollReqByDataKey(
          `${window.API_HOST}count/set_template/${this.templateId}/service_template/hosts`,
          { ids: serviceTemplateIds },
          { limit: 100 }
        )
        this.servicesHost = this.services.map(serviceTemplate => ({
          service_id: serviceTemplate.id,
          host_count: hostCounts.find(item => item.id === serviceTemplate.id)?.count
        }))
      },
      serviceExistHost(id) {
        const service = this.servicesHost.find(service => service.service_id === id)
        if (service) {
          return service.host_count > 0
        }
        return false
      },
      initMonitorTemplateName() {
        this.unwatch = this.$watch(() => this.$parent.templateName, (value) => {
          if (value) {
            this.templateName = value
          } else {
            this.templateName = this.$t('模板集群名称')
          }
        }, { immediate: true })
      },
      handleChildClick(service) {
        if (this.mode === 'view') {
          return false
        }
        this.selected = service.id
      },
      handleAddService() {
        this.selected = null
        this.dialog.props = {
          selected: this.services.map(service => service.id)
        }
        this.dialog.title = this.$t('添加服务模板')
        this.dialog.name = 'add'
        this.dialog.component = serviceTemplateSelector.name
        this.dialog.visible = true
      },
      handleViewService(service) {
        this.dialog.props = {
          id: service.id
        }
        this.dialog.title = `【${service.name}】${this.$t('模板服务信息')}`
        this.dialog.name = 'view'
        this.dialog.component = serviceTemplateInfo.name
        this.dialog.visible = true
      },
      handleDialogConfirm() {
        if (this.dialog.component === serviceTemplateSelector.name) {
          this.services = this.$refs.dialogComponent.getSelectedServices()
        }
        this.dialog.visible = false
      },
      handleDialogClose() {
        this.dialog.component = null
        this.dialog.title = ''
        this.dialog.props = {}
      },
      handleDeleteService(index) {
        this.services.splice(index, 1)
      },
      recoveryService() {
        this.services = [...this.originalServices]
      },
      handleGoTopoBusiness(service) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE,
          query: {
            keyword: service.name
          },
          history: true
        })
      },
      handleSelectChange(selected) {
        this.selectedServiceCount = selected.length
      },
      handleServiceTemplateLoaded(templates) {
        this.serviceTemplateCount = (templates || []).length
      },
      handleLinkClick() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    $iconColor: #C4C6CC;
    $fontColor: #63656E;
    $highlightColor: #3A84FF;
    $iconDisabledColor: #D8D8D8;
    .template-tree {
        padding: 10px 0 10px 20px;
        border: 1px dashed #C4C6CC;
        background-color: #FAFBFD;
        max-width: 960px;
        &:not(.view) {
            .node-child:hover {
                background-color: rgba(240, 241, 245, .6);
                .child-name {
                    color: $highlightColor;
                }
                .child-options {
                    display: block;
                }
            }
        }
    }
    .node-icon {
        position: relative;
        margin: 8px 4px 8px 0px;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        line-height: 20px;
        text-align: center;
        font-size: 12px;
        font-style: normal;
        color: #fff;
        background-color: #97AED6;
        z-index: 2;
    }
    .node-root {
        line-height: 36px;
        cursor: default;
        .folder-icon {
            width: 23px;
            height: 36px;
            line-height: 36px;
            text-align: center;
            font-size: 12px;
            color: $iconColor;
        }
        .root-name {
            display: block;
            padding: 0 10px 0 0;
            font-size: 14px;
            color: $fontColor;
            @include ellipsis;
        }
    }
    .node-children {
        line-height: 36px;
        margin-left: 32px;
        cursor: default;
        .node-child {
            padding: 0 10px 0 32px;
            position: relative;
            &.selected {
                background-color: #F0F1F5;
            }
            &.selected {
                .node-icon {
                    background-color: $highlightColor;
                }
                .child-name {
                    color: $highlightColor;
                }
                .child-options {
                    display: block;
                }
            }
            &:before {
                position: absolute;
                left: 0px;
                top: -18px;
                content: "";
                width: 25px;
                height: 36px;
                border-left: 1px dashed #DCDEE5;
                border-bottom: 1px dashed #DCDEE5;
                z-index: 1;
            }
            &.options-child {
                cursor: pointer;
                .node-icon {
                    font-size: 18px;
                    background-color: transparent;
                    color: $highlightColor;
                }
                .child-name {
                    color: $highlightColor;
                }
            }
            .child-name {
                display: block;
                padding: 0 10px 0 0;
                font-size: 14px;
                color: $fontColor;
                @include ellipsis;
            }
            .child-options {
                display: none;
                margin-right: 9px;
                font-size: 0;
                color: $iconColor;
                .action-link {
                  margin: 0 6px;
                  ::v-deep .bk-link-text {
                    font-size: 12px;
                  }

                  &:not(.disabled):hover {
                    color: $highlightColor;
                  }
                  &.disabled:hover {
                    cursor: not-allowed;
                  }
                }
            }
        }
    }
    .service-tips {
        span {
            color: $highlightColor;
            cursor: pointer;
        }
    }
    .dialog-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .summary {
            display: flex;
            font-size: 14px;
            .num {
                color: #2DCB56;
                font-style: normal;
                font-weight: 700;
                margin: 0 3px;
            }
            .to-template {
                margin-left: 16px;
                /deep/ .bk-link-text {
                    font-size: 12px;
                }
                /deep/ .bk-link-icon {
                    font-size: 12px;
                    margin-top: 2px;
                }
            }
        }
        .action {
            .btn {
                min-width: 76px;
            }
        }
    }
</style>
