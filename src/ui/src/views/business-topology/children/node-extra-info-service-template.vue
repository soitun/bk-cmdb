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
  <div class="clearfix" v-if="template">
    <div class="info-item fl" :title="`${$t('服务模板')} : ${template.name}`">
      <span class="name fl">{{$t('服务模板')}}</span>
      <div class="value fl">
        <div class="template-value" v-if="instance.service_template_id" @click="goServiceTemplate">
          <span class="text link">{{template.name}}</span>
          <i class="icon-cc-share"></i>
        </div>
        <span class="text" v-else>{{template.name}}</span>
      </div>
    </div>
    <div class="info-item fl" :title="`${$t('服务分类')} : ${template.category || '--'}`">
      <span class="name fl">{{$t('服务分类')}}</span>
      <div class="value fl">
        <span class="text">{{template.category || '--'}}</span>
      </div>
    </div>
    <div class="info-item fl" v-if="!isBizSet">
      <span class="name fl">{{$t('主机属性自动应用')}}</span>
      <div class="fl">
        <div class="template-value" @click="linkToAutoApply">
          <span class="text link">{{autoApplyEnable ? $t('已启用') : $t('未启用')}}</span>
          <i class="icon-cc-share link"></i>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { MENU_BUSINESS_HOST_APPLY, MENU_BUSINESS_SET_TOPOLOGY, MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'
  const serviceCategoryRequestId = Symbol('serviceCategoryRequestId')
  export default {
    name: 'service-template-info',
    props: {
      instance: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        template: null,
        serviceTemplateHostApplyEnabled: null
      }
    },
    computed: {
      isBizSet() {
        return this.$route.name === MENU_BUSINESS_SET_TOPOLOGY
      },
      bizId() {
        const { objectBiz, bizSet } = this.$store.state
        return this.isBizSet ? bizSet.bizId : objectBiz.bizId
      },
      moduleId() {
        return this.selectedNode?.data?.bk_inst_id
      },
      selectedNode() {
        return this.$store.state.businessHost.selectedNode
      },
      autoApplyEnable() {
        // 在不同版本的拓扑数据接口中，node节点中是否存在主机自动应用字段是不一致的
        // 所以，统一通过instance获取，保证数据获取的正确性
        return this.instance.host_apply_enabled || this.serviceTemplateHostApplyEnabled
      }
    },
    watch: {
      instance: {
        immediate: true,
        handler(instance) {
          instance && this.setInfo()
        }
      }
    },
    methods: {
      setInfo() {
        this.getServiceInfo()

        // 使用了服务模板需要查看模块所属的服务模板是否开启了主机自动应用
        if (this.instance.service_template_id && !this.isBizSet) {
          this.getModuleApplyStatusByTemplate()
        }
      },
      async getServiceInfo() {
        const categories = await this.getServiceCategories()
        // eslint-disable-next-line max-len
        const firstCategory = categories.find(category => category.secondCategory.some(second => second.id === this.instance.service_category_id)) || {}
        // eslint-disable-next-line max-len
        const secondCategory = (firstCategory.secondCategory || []).find(second => second.id === this.instance.service_category_id) || {}
        this.template = {
          name: this.instance.service_template_id ? this.instance.bk_module_name : this.$t('无'),
          category: `${firstCategory.name || '--'} / ${secondCategory.name || '--'}`
        }
      },
      async getModuleApplyStatusByTemplate() {
        const [moduleApplyStatus] = await this.$store.dispatch('hostApply/getModuleApplyStatusByTemplate', {
          params: {
            bk_biz_id: this.bizId,
            bk_module_ids: [this.moduleId]
          }
        })

        this.serviceTemplateHostApplyEnabled = moduleApplyStatus?.host_apply_enabled
      },
      async getServiceCategories() {
        try {
          const { info = [] } = await this.$store.dispatch('serviceClassification/searchServiceCategory', {
            params: { bk_biz_id: this.bizId },
            config: {
              requestId: serviceCategoryRequestId,
              fromCache: true
            }
          })
          const categories = this.collectServiceCategories(info)
          return categories
        } catch (e) {
          console.error(e)
          return []
        }
      },
      collectServiceCategories(data) {
        const categories = []
        data.forEach((item) => {
          if (!item.category.bk_parent_id) {
            categories.push(item.category)
          }
        })
        categories.forEach((category) => {
          // eslint-disable-next-line max-len
          category.secondCategory = data.filter(item => item.category.bk_parent_id === category.id).map(item => item.category)
        })
        return categories
      },
      async goServiceTemplate() {
        try {
          const data = await this.$store.dispatch('serviceTemplate/findServiceTemplate', {
            id: this.instance.service_template_id,
            config: {
              globalError: false
            }
          })
          if (!data) {
            return this.$error(this.$t('跳转失败，服务模板已经被删除'))
          }

          this.$routerActions.open({
            name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
            params: {
              templateId: this.instance.service_template_id,
              bizId: data?.template?.bk_biz_id,
            }
          })
        } catch (error) {
          console.error(error)
          this.$error(error.message)
        }
      },
      linkToAutoApply() {
        this.$routerActions.open({
          name: MENU_BUSINESS_HOST_APPLY,
          params: {
            bizId: this.bizId,
            mode: 'module'
          },
          query: {
            id: this.moduleId
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .info-item {
        width: 50%;
        max-width: 400px;
        line-height: 26px;
        margin-bottom: 12px;
        .name {
            position: relative;
            padding-right: 16px;
            &::after {
                content: ":";
                position: absolute;
                right: 10px;
            }
        }
        .value {
            width: calc(80% - 10px);
            padding-right: 10px;
            .text {
                @include inlineBlock;
                @include ellipsis;
                max-width: calc(100% - 16px);
                font-size: 14px;
            }
            .template-value {
                width: 100%;
                font-size: 0;
                color: #3a84ff;
                cursor: pointer;
            }
        }
        .icon-cc-share {
            @include inlineBlock;
            font-size: 12px;
            margin-left: 4px;
        }
        .link {
            color: #3a84ff;
            cursor: pointer;
        }
    }
</style>
