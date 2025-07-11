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
    <h2 class="node-create-title">{{title}}</h2>
    <div class="node-create-path">{{$t('添加节点已选择')}}：{{topoPath}}</div>
    <div class="node-create-form"
      :style="{
        'max-height': Math.min($APP.height - 400, 400) + 'px',
        'padding-bottom': formPaddingBottom
      }">
      <div v-for="(property, index) in sortedProperties"
        :class="[
          'form-group',
          { 'form-group-flex': sortedProperties.length === 1 || property['bk_property_type'] === 'longchar' }
        ]"
        :key="index">
        <label :class="['form-label', 'inline-block-middle', {
          required: property['isrequired']
        }]">
          {{property['bk_property_name']}}
        </label>
        <component v-if="!['longchar'].includes(property['bk_property_type'])"
          :is="`cmdb-form-${property['bk_property_type']}`"
          style="display: block;"
          :unit="property['unit']"
          :data-vv-name="property['bk_property_id']"
          :data-vv-as="property['bk_property_name']"
          :options="property.option || []"
          :placeholder="$t('请输入xx', { name: property.bk_property_name })"
          v-validate="getValidateRules(property)"
          v-bind="$tools.getValidateEvents(property)"
          v-model.trim="values[property['bk_property_id']]">
        </component>
        <div v-else>
          <bk-input type="textarea" class="longchar-textarea"
            :data-vv-name="property['bk_property_id']"
            :data-vv-as="property['bk_property_name']"
            :options="property.option || []"
            :placeholder="$t('请输入xx', { name: property.bk_property_name })"
            v-validate="getValidateRules(property)"
            v-bind="$tools.getValidateEvents(property)"
            v-model.trim="values[property['bk_property_id']]">
          </bk-input>
        </div>
        <span class="form-error">{{errors.first(property['bk_property_id'])}}</span>
      </div>
    </div>
    <div class="node-create-options">
      <bk-button theme="primary" v-test-id="'createNodeSave'"
        :disabled="$loading() || errors.any()"
        @click="handleSave">
        {{$t('提交')}}
      </bk-button>
      <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import { BUILTIN_UNEDITABLE_FIELDS } from '@/dictionary/model-constants'
  export default {
    props: {
      parentNode: {
        type: Object,
        required: true
      },
      properties: {
        type: Array,
        required: true
      },
      nextModelId: String
    },
    data() {
      return {
        values: {}
      }
    },
    computed: {
      topoPath() {
        const nodePath = [...this.parentNode.parents, this.parentNode]
        return nodePath.map(node => node.data.bk_inst_name).join('-')
      },
      sortedProperties() {
        return this.properties.filter(property => !BUILTIN_UNEDITABLE_FIELDS.includes(property.bk_property_id))
          .sort((propertyA, propertyB) => propertyA.bk_property_index - propertyB.bk_property_index)
      },
      title() {
        return this.nextModelId === 'set' ? this.$t('新建集群') : this.$t('新建节点')
      },
      formPaddingBottom() {
        return this.nextModelId === 'set' ? '20px' : '52px'
      }
    },
    watch: {
      properties() {
        this.initValues()
      }
    },
    created() {
      this.initValues()
    },
    methods: {
      initValues() {
        this.values = this.$tools.getInstFormValues(this.properties, {})
      },
      getValidateRules(property) {
        const rules = this.$tools.getValidateRules(property)
        if (property.bk_property_id === 'bk_inst_name') {
          rules.businessTopoInstNames = true
          rules.length = 256
          rules.singlechar = false
        }
        return rules
      },
      handleSave() {
        this.$validator.validateAll().then((isValid) => {
          if (isValid) {
            this.$emit('submit', this.$tools.formatValues(this.values, this.properties))
          }
        })
      },
      handleCancel() {
        this.$emit('cancel')
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
        padding: 12px 26px 18px;
        font-size: 12px;
    }
    .node-create-form {
        padding: 0 26px;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        @include scrollbar-y;
        .form-group {
            flex: 0 0 48%;
            position: relative;
            padding: 0 0 16px;
            .form-label {
                position: relative;
                display: inline-block;
                padding: 0 10px 0 0;
                max-width: 100%;
                line-height: 24px;
                @include ellipsis;
                &.required:after {
                    position: absolute;
                    right: 0;
                    content: '*';
                    color: #ff5656;
                }
            }
            .longchar-textarea {
                /deep/ .bk-form-textarea {
                    min-height: 60px !important;
                    height: 60px !important;
                }
            }
            .form-error {
                position: absolute;
                bottom: -2px;
                left: 0;
                font-size: 12px;
                color: #ff5656;
            }
            &.form-group-flex {
                flex: 0 0 100%;
            }
        }
    }
    .node-create-options {
        padding: 9px 20px;
        border-top: 1px solid $cmdbBorderColor;
        text-align: right;
        background-color: #fafbfd;
    }
</style>
