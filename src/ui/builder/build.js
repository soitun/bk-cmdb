/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

process.env.NODE_ENV = 'production'

const webpack = require('webpack')
const webpackConfig = require('./webpack')

webpack(webpackConfig, (err, stats) => {
  if (err) {
    console.error(err)
    return
  }

  if (stats.hasErrors()) {
    stats.compilation.errors.forEach((e) => {
      console.error(e.message)
    })
    return
  }

  console.log(stats.toString({ colors: true }))
})
