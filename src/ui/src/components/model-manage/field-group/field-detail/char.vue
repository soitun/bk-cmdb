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
  <div class="form-label cmdb-form-item" :class="{ 'is-error': errors.has('option') }">
    <span class="label-text">{{$t('正则校验')}}</span>
    <bk-input
      type="textarea"
      class="raw"
      name="option"
      v-model="localValue"
      :disabled="isReadOnly"
      data-vv-validate-on="blur"
      v-validate="'remoteRegular'"
      @input="handleInput">
    </bk-input>
    <p class="form-error">{{errors.first('option')}}</p>
  </div>
</template>

<script>
  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      isReadOnly: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        localValue: ''
      }
    },
    watch: {
      value() {
        this.localValue = this.value === '' ? '' : this.value
      }
    },
    created() {
      this.localValue = this.value === '' ? '' : this.value
    },
    methods: {
      handleInput() {
        this.$emit('input', this.localValue)
      },
      validate() {
        return this.$validator.validateAll()
      }
    }
  }
</script>
