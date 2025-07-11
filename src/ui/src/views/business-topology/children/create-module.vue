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
  <div class="node-create-layout">
    <h2 class="node-create-title">{{$t('新建模块')}}</h2>
    <div class="node-create-path" :title="topoPath">{{$t('添加节点已选择')}}：{{topoPath}}</div>
    <div class="node-create-form"
      :style="{
        'max-height': Math.min($APP.height - 400, 400) + 'px'
      }">
      <div class="form-item clearfix mt30">
        <div class="create-type fl">
          <input class="type-radio"
            type="radio"
            id="formTemplate"
            name="createType"
            v-model="withTemplate"
            :value="1">
          <label for="formTemplate">{{$t('从模板新建')}}</label>
        </div>
        <div class="create-type fl ml50">
          <input class="type-radio"
            type="radio"
            id="createDirectly"
            name="createType"
            v-model="withTemplate"
            :value="0">
          <label for="createDirectly">{{$t('直接新建')}}</label>
        </div>
      </div>
      <div class="form-item" v-if="withTemplate">
        <label>{{$t('服务模板')}}</label>
        <bk-select style="width: 100%;"
          :clearable="false"
          :searchable="templateList.length > 7"
          :loading="$loading(request.serviceTemplate)"
          v-model="template"
          v-validate="'required'"
          data-vv-name="template"
          key="template">
          <bk-option v-for="(option, index) in templateList"
            :key="index"
            :id="option.id"
            :name="option.name">
          </bk-option>
          <div class="add-template" slot="extension" @click="jumpServiceTemplate" v-if="!templateList.length">
            <i class="bk-icon icon-plus-circle"></i>
            <span>{{$t('新建服务模板')}}</span>
          </div>
        </bk-select>
        <span class="form-error" v-if="errors.has('template')">{{errors.first('template')}}</span>
      </div>
      <div class="form-item">
        <label>
          {{$t('模块名称')}}
          <font color="red">*</font>
          <i class="icon-cc-tips"
            v-bk-tooltips.top="$t('模块名称提示')"
            v-if="withTemplate === 1">
          </i>
        </label>
        <cmdb-form-singlechar
          v-if="withTemplate"
          v-model="moduleName"
          v-validate="'required|businessTopoInstNames|length:256'"
          data-vv-name="moduleName"
          data-vv-validate-on="blur"
          key="moduleName"
          :placeholder="$t('请输入xx', { name: $t('模块名称') })"
          :disabled="!!withTemplate">
        </cmdb-form-singlechar>
        <bk-input class="form-textarea" v-else
          type="textarea"
          data-vv-name="moduleNameMulti"
          v-validate="
            `required|longchar|businessTopoInstNames|emptyModuleName|moduleNameMap|moduleNameLen|splitMaxLength:100
            ,${$t('超过限制，一次最多支持创建n个', { n: 100 })}
          `"
          v-model="moduleNameMulti"
          :rows="rows"
          :placeholder="!!withTemplate ? $t('请输入xx', { name: $t('模块名称') }) : $t('模块多个创建提示')"
          :disabled="!!withTemplate"
          @keydown="handleKeydown"
          @paste="handlePaste">
        </bk-input>
        <span class="form-error" v-if="errors.has('moduleName')">{{errors.first('moduleName')}}</span>
        <span class="form-error" v-if="errors.has('moduleNameMulti')">{{errors.first('moduleNameMulti')}}</span>
      </div>
      <div class="form-item clearfix" v-if="!withTemplate">
        <label>{{$t('所属服务分类')}}<font color="red">*</font></label>
        <cmdb-selector class="service-class fl"
          v-model="firstClass"
          v-validate="'required'"
          data-vv-name="firstClass"
          key="firstClass"
          :auto-select="false"
          :list="firstClassList"
          :loading="$loading(request.serviceCategory)"
          @on-selected="updateCategory">
        </cmdb-selector>
        <cmdb-selector class="service-class fr"
          v-model="secondClass"
          v-validate="'required'"
          data-vv-name="secondClass"
          key="secondClass"
          :list="secondClassList"
          :loading="$loading(request.serviceCategory)">
        </cmdb-selector>
        <span class="form-error" v-if="errors.has('firstClass')">{{errors.first('firstClass')}}</span>
        <span class="form-error second-class" v-if="errors.has('secondClass')">{{errors.first('secondClass')}}</span>
      </div>
    </div>
    <div class="node-create-options">
      <bk-button theme="primary" v-test-id="'createModuleSave'"
        :disabled="$loading() || errors.any()"
        @click="handleSave">
        {{$t('提交')}}
      </bk-button>
      <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import has from 'has'
  import serviceTemplateService from '@/service/service-template/index.js'
  import { MENU_BUSINESS_SERVICE_TEMPLATE } from '@/dictionary/menu-symbol'

  export default {
    props: {
      parentNode: {
        type: Object,
        required: true
      }
    },
    data() {
      return {
        withTemplate: 1,
        createTypeList: [{
          id: 1,
          name: this.$t('从模板创建')
        }, {
          id: 0,
          name: this.$t('直接创建')
        }],
        template: '',
        templateList: [],
        moduleName: '',
        moduleNameMulti: '',
        rows: 1,
        firstClass: '',
        firstClassList: [],
        secondClass: '',
        values: {},
        request: {
          serviceTemplate: Symbol('serviceTemplate'),
          serviceCategory: Symbol('serviceCategory')
        }
      }
    },
    computed: {
      topoPath() {
        const nodePath = [...this.parentNode.parents, this.parentNode]
        return nodePath.map(node => node.data.bk_inst_name).join('/')
      },
      business() {
        return this.$store.getters['objectBiz/bizId']
      },
      serviceTemplateMap() {
        return this.$store.state.businessHost.serviceTemplateMap
      },
      currentTemplate() {
        return this.templateList.find(item => item.id === this.template) || {}
      },
      categoryMap() {
        return this.$store.state.businessHost.categoryMap
      },
      currentCategory() {
        return this.firstClassList.find(category => category.id === this.firstClass) || {}
      },
      secondClassList() {
        return this.currentCategory.secondCategory || []
      }
    },
    watch: {
      withTemplate(withTemplate) {
        if (withTemplate) {
          this.updateCategory()
          this.template = this.templateList.length ? this.templateList[0].id : ''
        } else {
          this.template = ''
          this.updateCategory(1)
          this.getServiceCategories()
        }
      },
      template(template) {
        if (template) {
          this.moduleName = this.currentTemplate.name
        } else {
          this.moduleName = ''
        }
      }
    },
    created() {
      this.getServiceTemplates()
    },
    methods: {
      setRows() {
        setTimeout(() => {
          const rows = this.moduleNameMulti.split('\n').length
          this.rows = Math.min(3, Math.max(rows, 1))
        })
      },
      handleKeydown(value, keyEvent) {
        if (['Enter', 'NumpadEnter'].includes(keyEvent.code)) {
          this.rows = Math.min(this.rows + 1, 3)
        } else if (keyEvent.code === 'Backspace') {
          this.setRows()
        }
      },
      handlePaste() {
        this.setRows()
      },
      async getServiceTemplates() {
        if (has(this.serviceTemplateMap, this.business)) {
          this.templateList = this.serviceTemplateMap[this.business]
        } else {
          try {
            const templates = await serviceTemplateService.findAll({ bk_biz_id: this.business, page: { sort: '-last_time' } }, {
              requestId: this.request.serviceTemplate
            })

            this.templateList = templates
            this.$store.commit('businessHost/setServiceTemplate', {
              id: this.business,
              templates
            })
          } catch (e) {
            console.error(e)
            this.templateList = []
          }
        }
        this.template = this.templateList[0] ? this.templateList[0].id : ''
      },
      async getServiceCategories() {
        if (has(this.categoryMap, this.business)) {
          this.firstClassList = this.categoryMap[this.business]
        } else {
          try {
            const data = await this.$store.dispatch('serviceClassification/searchServiceCategory', {
              params: { bk_biz_id: this.business },
              config: {
                requestId: this.request.serviceCategory
              }
            })
            const categories = this.collectServiceCategories(data.info)
            this.firstClassList = categories
            this.$store.commit('businessHost/setCategories', {
              id: this.business,
              categories
            })
          } catch (e) {
            console.error(e)
            this.firstClassList = []
          }
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
      updateCategory(firstClass) {
        if (firstClass) {
          this.firstClass = firstClass
          this.secondClass = this.secondClassList.length ? this.secondClassList[0].id : ''
        } else {
          this.firstClass = ''
          this.secondClass = ''
        }
      },
      handleSave() {
        this.$validator.validateAll().then((isValid) => {
          if (isValid) {
            const data = {
              service_category_id: this.withTemplate ? this.currentTemplate.service_category_id : this.secondClass,
              service_template_id: this.withTemplate ? this.template : 0
            }
            if (this.withTemplate) {
              data.bk_module_name = this.moduleName
            } else {
              const nameList = this.moduleNameMulti.split('\n').filter(name => name.trim().length)
                .map(name => name.trim())
              data.bk_module_name = nameList
            }
            this.$emit('submit', data)
          }
        })
      },
      handleCancel() {
        this.$emit('cancel')
      },
      jumpServiceTemplate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE,
          params: {
            bizId: this.business
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .node-create-layout {
        position: relative;
    }
    .node-create-title {
        margin-top: -14px;
        padding: 0 26px;
        line-height: 30px;
        font-size: 22px;
        color: #333948;
    }
    .node-create-path {
        padding: 23px 26px 0;
        margin: 0 0 -5px 0;
        font-size: 12px;
        @include ellipsis;
    }
    .node-create-form {
        padding: 0 26px 27px;
        overflow: visible;
    }
    .form-item {
        margin: 15px 0 0 0;
        position: relative;
        label {
            display: block;
            padding: 7px 0;
            line-height: 19px;
            font-size: 14px;
            color: #63656E;
        }
        .service-class {
            width: 260px;
            @include inlineBlock;
        }
        .form-textarea {
            /deep/ textarea {
                min-height: auto !important;
                line-height: 22px;
                @include scrollbar-y(6px);
            }
        }
        .form-error {
            position: absolute;
            top: 100%;
            left: 0;
            font-size: 12px;
            color: $cmdbDangerColor;
            &.second-class {
                left: 270px;
            }
        }
        .create-type {
            display: flex;
            align-items: center;
            .type-radio {
                -webkit-appearance: none;
                width: 16px;
                height: 16px;
                padding: 3px;
                border: 1px solid #979BA5;
                border-radius: 50%;
                background-clip: content-box;
                outline: none;
                cursor: pointer;
                &:checked {
                    border-color: #3A84FF;
                    background-color: #3A84FF;
                }
            }
            label {
                padding: 0 0 0 6px;
                font-size: 14px;
                cursor: pointer;
            }
        }
    }
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #fafbfd;
    }
    font {
        padding: 0 2px;
    }
    .add-template {
        width: 20%;
        cursor: pointer;
        .icon-plus-circle {
            @include inlineBlock;
            font-size: 14px;
        }
        span {
            @include inlineBlock;
        }
    }
</style>
