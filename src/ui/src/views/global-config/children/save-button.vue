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
  <bk-popconfirm trigger="click" :disabled="loading || disabled" :title="$t('确认保存配置？')"
    :content="$t('保存后将会立即生效，请谨慎操作！')" @confirm="handleConfirm">
    <bk-button theme="primary" class="action-button" :loading="loading" :disabled="disabled">{{$t('保存')}}</bk-button>
  </bk-popconfirm>
</template>

<script>
  export default {
    props: {
      loading: {
        type: Boolean,
        default: false
      },
      disabled: {
        type: Boolean,
        default: false
      }
    },
    setup(props, { emit }) {
      const handleConfirm = () => emit('save')

      return { handleConfirm }
    }
  }
</script>

<style lang="scss" scoped>
@import url("../style.scss");
</style>
