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
  <div class="verification-layout">
    <div class="options">
      <cmdb-auth class="inline-block-middle"
        v-if="!isTopoModel && isShowOptionBtn"
        :auth="{ type: $OPERATION.U_MODEL, relation: [modelId] }"
        @update-auth="handleReceiveAuth">
        <bk-button slot-scope="{ disabled }"
          class="create-btn"
          theme="primary"
          :disabled="isReadOnly || disabled"
          @click="createVerification">
          {{$t('新建校验')}}
        </bk-button>
      </cmdb-auth>
    </div>
    <bk-table
      class="verification-table"
      v-bkloading="{
        isLoading: $loading(['searchObjectUniqueConstraints', 'deleteObjectUniqueConstraints'])
      }"
      :data="table.list"
      :max-height="$APP.height - 320"
      :row-style="{
        cursor: 'pointer'
      }"
      @cell-click="handleShowDetails">
      <bk-table-column :label="$t('校验规则')" class-name="is-highlight" show-overflow-tooltip>
        <template slot-scope="{ row }">
          <div class="keys-cell">
            {{getRuleName(row.keys)}}
            <mini-tag :text="$t('模板')" v-if="row.bk_template_id" />
          </div>
        </template>
      </bk-table-column>
      <bk-table-column prop="operation"
        v-if="updateAuth && !isTopoModel"
        :label="$t('操作')">
        <template slot-scope="{ row }">
          <button class="text-primary mr10 operation-btn"
            :disabled="!isEditable(row)"
            @click.stop="editVerification(row)">
            {{$t('编辑')}}
          </button>
          <button class="text-primary operation-btn"
            :disabled="!isEditable(row)"
            @click.stop="deleteVerification(row)">
            {{$t('删除')}}
          </button>
        </template>
      </bk-table-column>
      <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
    </bk-table>
    <bk-sideslider
      v-transfer-dom
      :width="450"
      :title="slider.title"
      :is-show.sync="slider.isShow"
      :before-close="handleSliderBeforeClose">
      <the-verification-detail
        ref="verificationForm"
        slot="content"
        v-if="slider.isShow"
        :attribute-list="attributeList"
        :rule-list="table.list"
        :id="slider.id"
        :readonly="slider.readonly"
        @save="saveVerification"
        @cancel="handleSliderBeforeClose"
        @change-mode="editVerification">
      </the-verification-detail>
    </bk-sideslider>
  </div>
</template>

<script>
  import theVerificationDetail from './verification-detail'
  import { mapActions, mapGetters } from 'vuex'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import MiniTag from '@/components/ui/other/mini-tag.vue'

  export default {
    components: {
      theVerificationDetail,
      MiniTag
    },
    props: {
      modelId: {
        type: Number,
        default: null
      }
    },
    data() {
      return {
        slider: {
          isShow: false,
          id: null,
          readonly: false
        },
        table: {
          list: [],
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        attributeList: [],
        updateAuth: false
      }
    },
    computed: {
      ...mapGetters('objectModel', [
        'activeModel'
      ]),
      ...mapGetters('objectMainLineModule', ['isMainLine']),
      isTopoModel() {
        // 主线模型除主机外
        return this.isMainLine(this.activeModel) && this.activeModel.bk_obj_id !== BUILTIN_MODELS.HOST
      },
      isReadOnly() {
        if (this.activeModel) {
          return this.activeModel.bk_ispaused
        }
        return false
      },
      isShowOptionBtn() {
        return BUILTIN_MODELS.PROJECT !== this.$route.params.modelId
      }
    },
    watch: {
      activeModel: {
        immediate: true,
        async handler(activeModel) {
          if (activeModel.bk_obj_id) {
            await this.initAttrList()
            this.searchVerification()
          }
        }
      }
    },
    methods: {
      ...mapActions('objectModelProperty', [
        'searchObjectAttribute'
      ]),
      ...mapActions('objectUnique', [
        'searchObjectUniqueConstraints',
        'deleteObjectUniqueConstraints'
      ]),
      isEditable(item) {
        if (item.ispre || this.isReadOnly || item.bk_template_id) {
          return false
        }
        return true
      },
      getRuleName(keys) {
        const name = []
        keys.forEach((key) => {
          if (key.key_kind === 'property') {
            const attr = this.attributeList.find(({ id }) => id === key.key_id)
            if (attr) {
              name.push(attr.bk_property_name)
            }
          }
        })
        return name.join('+')
      },
      async initAttrList() {
        this.attributeList = await this.searchObjectAttribute({
          params: {
            bk_obj_id: this.activeModel.bk_obj_id
          },
          config: {
            requestId: `post_searchObjectAttribute_${this.activeModel.bk_obj_id}`
          }
        })
      },
      createVerification() {
        this.slider.title = this.$t('新建校验')
        this.slider.id = null
        this.slider.readonly = false
        this.slider.isShow = true
      },
      editVerification({ id }) {
        this.slider.title = this.$t('编辑校验')
        this.slider.id = id
        this.slider.readonly = false
        this.slider.isShow = true
      },
      saveVerification() {
        this.slider.isShow = false
        this.slider.id = null
        this.slider.readonly = false
        this.searchVerification()
      },
      deleteVerification(verification) {
        this.$bkInfo({
          title: this.$tc('确定删除唯一校验', this.getRuleName(verification.keys), { name: this.getRuleName(verification.keys) }),
          confirmFn: async () => {
            await this.deleteObjectUniqueConstraints({
              objId: verification.bk_obj_id,
              id: verification.id,
              params: {},
              config: {
                requestId: 'deleteObjectUniqueConstraints'
              }
            })
            this.searchVerification()
          }
        })
      },
      async searchVerification() {
        const uniqueList = await this.searchObjectUniqueConstraints({
          objId: this.activeModel.bk_obj_id,
          params: {},
          config: {
            requestId: 'searchObjectUniqueConstraints'
          }
        })

        // 只保留在对象属性列表中能找到的规则字段，如果某条记录一个字段都找不到则不会显示
        const list = uniqueList
          .filter(item => item.keys.every(key => this.attributeList.find(({ id }) => id === key.key_id)))

        this.table.list = list
      },
      handleShowDetails(row, column) {
        if (column.property === 'operation') return
        this.slider.title = this.$t('查看校验')
        this.slider.id = row.id
        this.slider.readonly = true
        this.slider.isShow = true
      },
      handleReceiveAuth(auth) {
        this.updateAuth = auth
      },
      handleSliderBeforeClose() {
        const isChanged = this.$refs.verificationForm.isChanged()
        const confirmFn = () => {
          this.slider.isShow = false
          this.slider.id = null
        }
        if (!isChanged) {
          confirmFn()
          return true
        }
        // this.$refs.verificationForm.setChanged(true)
        return this.$refs.verificationForm.beforeClose(confirmFn)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .verification-layout {
        padding: 20px;
    }
    .verification-table {
        margin: 14px 0 0 0;
        .keys-cell {
          display: flex;
          align-items: center;
          gap: 4px;
        }
    }
    .operation-btn[disabled] {
        color: #dcdee5 !important;
        opacity: 1 !important;
    }
</style>
