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

<template functional>
  <div class="no-search-results">
    <img src="../../assets/images/full-text-search.png" alt="no-search-results">
    <p>{{props.text}}</p>
  </div>
</template>

<style lang="scss" scoped>
  .no-search-results {
    width: 90%;
    margin: 0 auto;
    padding: 32px 0;
    text-align: center;
    color: #63656E;
    font-size: 16px;
    img {
      width: 104px;
    }
  }
</style>
