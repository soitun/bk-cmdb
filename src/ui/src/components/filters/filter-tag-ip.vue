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
  <span class="filter-tag" @click="handleClick"
    v-bk-tooltips="{
      theme: 'filter-tag',
      disabled: value.length < 3,
      content: value.join('<br>'),
      interactive: true,
      hideOnClick: false,
      allowHTML: true,
      onShow(inst) {
        inst.setContent(value.join('<br>'))
      }
    }"
  >
    <label class="tag-name">{{label}}</label>
    <span class="tag-colon">:</span>
    <span class="tag-value">{{displayText}}</span>
    <i class="tag-delete bk-icon icon-close" @mouseenter.prevent.stop @click.stop="handleRemove"></i>
  </span>
</template>

<script>
  import FilterStore from './store'
  import Utils from './utils'
  import FilterForm from './filter-form.js'
  export default {
    computed: {
      label() {
        const { inner, outer, exact } = FilterStore.IP
        const label = []
        inner && label.push(this.$t('内网IP'))
        outer && label.push(this.$t('外网IP'))
        exact && label.push(this.$t('精确'))
        return label.join(' | ')
      },
      value() {
        return Utils.splitIP(FilterStore.IP.text)
      },
      displayText() {
        const count = this.value.length
        const [firstIp] = this.value
        if (count > 2) {
          return this.$i18n.locale === 'en' ? `${firstIp}...etc ${count} IPs in all` : `${firstIp}...等${count}个IP`
        }
        return this.value.join(' | ')
      }
    },
    methods: {
      handleClick() {
        FilterForm.show()
      },
      handleRemove() {
        FilterStore.resetIP()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .filter-tag {
        display: inline-flex;
        align-items: center;
        margin: 0 3px 10px;
        padding: 0 0 0 5px;
        border-radius: 2px;
        font-size: 12px;
        background: #f0f1f5;
        line-height: 22px;
        cursor: pointer;
        &:hover {
            background-color: #DCDEE5;
        }
        .tag-name {
            max-width: 150px;
            padding-right: 5px;
            color: #63656E;
            cursor: pointer;
            @include ellipsis;
        }
        .tag-colon {
            padding-right: 5px;
        }
        .tag-value {
            max-width: 255px;
            color: #313238;
            @include ellipsis;
        }
        .tag-delete {
            font-size: 20px;
            color: #9b9ea8;
            &:hover {
                color: #313238;
            }
        }
    }
</style>
<style lang="scss">
.tippy-tooltip.filter-tag-theme {
  padding: 7px;
  .tippy-content {
    max-height: 300px;
    padding: 0 8px;
    @include scrollbar-y(6px);
  }
}
</style>
