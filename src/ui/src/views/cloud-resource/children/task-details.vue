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
  <div v-bkloading="{ isLoading: $loading(request.findOneTask) }">
    <bk-tab class="details-tab" :active.sync="active" type="unborder-card" slot="content">
      <bk-tab-panel name="details" :label="$t('任务详情')">
      </bk-tab-panel>
      <bk-tab-panel name="history" :label="$t('录入历史')">
      </bk-tab-panel>
    </bk-tab>
    <component class="details-component" :is="component"
      ref="component"
      :container="this"
      :task="task"
      :id="id">
    </component>
  </div>
</template>

<script>
  import TaskDetailsInfo from './task-details-info.vue'
  import TaskDetailsHistory from './task-details-history.vue'
  import TaskForm from './task-form.vue'
  import useSideslider from '@/hooks/use-sideslider'

  export default {
    name: 'task-details',
    components: {
      [TaskDetailsHistory.name]: TaskDetailsHistory,
      [TaskDetailsInfo.name]: TaskDetailsInfo,
      [TaskForm.name]: TaskForm
    },
    props: {
      id: {
        type: Number,
        required: true
      },
      defaultComponent: String,
      container: Object
    },
    data() {
      return {
        task: null,
        title: '',
        active: 'details',
        detailsComponent: this.defaultComponent || TaskDetailsInfo.name,
        request: {
          findOneTask: Symbol('findOneTask')
        }
      }
    },
    computed: {
      component() {
        if (!this.task) {
          return null
        }
        if (this.active === 'details') {
          return this.detailsComponent
        }
        return TaskDetailsHistory.name
      }
    },
    created() {
      this.getTaskDetails()
      const { beforeClose, setChanged  } = useSideslider()
      this.beforeClose = beforeClose
      this.setChanged = setChanged
    },
    methods: {
      async getTaskDetails() {
        try {
          this.task = await this.$store.dispatch('cloud/resource/findOneTask', {
            id: this.id,
            config: {
              requestId: this.request.findOneTask
            }
          })
        } catch (e) {
          console.error(e)
          this.task = null
        }
      },
      show(options) {
        this.title = options.title || this.title
        this.task = options.task || this.task
        this.detailsComponent = options.detailsComponent || TaskDetailsInfo.name
      },
      hide(eventType) {
        eventType && this.$emit(eventType)
        this.container.hide()
      },
      changedValues() {
        const { changedValues } = this.$refs.component
        if (changedValues) {
          return changedValues()
        }
        return null
      },
      changedVpcValues() {
        const { changedVpcValues } = this.$refs.component
        if (changedVpcValues) {
          return changedVpcValues()
        }
        return null
      }
    }
  }
</script>

<style lang="scss" scoped>
    .details-tab {
        height: auto;
        /deep/ {
            .bk-tab-header {
                padding: 0;
                margin: 0 24px;
            }
            .bk-tab-section {
                height: 0;
            }
        }
    }
    .details-component {
        height: calc(100% - 58px);
        @include scrollbar-y;
    }
</style>
