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

import { t } from '@/i18n'

export const QUERY_OPERATOR = Object.freeze({
  EQ: '$eq',
  NE: '$ne',
  IN: '$in',
  NIN: '$nin',
  LT: '$lt',
  GT: '$gt',
  LTE: '$lte',
  GTE: '$gte',
  // 前端构造的操作符，真实数据中会拆分数据为gte, lte向后台传递
  RANGE: '$range',
  NRANGE: '$nrange',

  // 在页面展示时like有不同的名称，接口中主机为$regex其它资源为contains，在前端选择组件中值都是$regex
  LIKE: '$regex',

  CONTAINS: '$contains',
  CONTAINS_CS: '$contains_s',
})

export const QUERY_OPERATOR_SYMBOL = {
  [QUERY_OPERATOR.EQ]: '=',
  [QUERY_OPERATOR.NE]: '≠',
  [QUERY_OPERATOR.IN]: 'in',
  [QUERY_OPERATOR.NIN]: 'not in',
  [QUERY_OPERATOR.GT]: '>',
  [QUERY_OPERATOR.LT]: '<',
  [QUERY_OPERATOR.GTE]: '≥',
  [QUERY_OPERATOR.LTE]: '≤',
  [QUERY_OPERATOR.LIKE]: 'like',
  [QUERY_OPERATOR.RANGE]: '≤ ≥',
  [QUERY_OPERATOR.CONTAINS]: 'contains',
  [QUERY_OPERATOR.CONTAINS_CS]: 'contains(CS)',
}

export const QUERY_OPERATOR_HOST_SYMBOL = {
  ...QUERY_OPERATOR_SYMBOL,
  [QUERY_OPERATOR.LIKE]: QUERY_OPERATOR_SYMBOL[QUERY_OPERATOR.CONTAINS_CS],
}

export const QUERY_OPERATOR_OTHER_SYMBOL = {
  ...QUERY_OPERATOR_SYMBOL,
  [QUERY_OPERATOR.LIKE]: QUERY_OPERATOR_SYMBOL[QUERY_OPERATOR.CONTAINS],
}

export const QUERY_OPERATOR_DESC = {
  [QUERY_OPERATOR.EQ]: t('等于'),
  [QUERY_OPERATOR.NE]: t('不等于'),
  [QUERY_OPERATOR.LT]: t('小于'),
  [QUERY_OPERATOR.GT]: t('大于'),
  [QUERY_OPERATOR.IN]: t('包含在'),
  [QUERY_OPERATOR.NIN]: t('不包含在'),
  [QUERY_OPERATOR.RANGE]: t('数值范围'),
  [QUERY_OPERATOR.LTE]: t('小于等于'),
  [QUERY_OPERATOR.GTE]: t('大于等于'),
  [QUERY_OPERATOR.LIKE]: t('模糊'),
  [QUERY_OPERATOR.CONTAINS]: t('匹配(大小写不敏感)'),
  [QUERY_OPERATOR.CONTAINS_CS]: t('匹配(大小写敏感)')
}

export const QUERY_OPERATOR_HOST_DESC = {
  ...QUERY_OPERATOR_DESC,
  [QUERY_OPERATOR.LIKE]: QUERY_OPERATOR_DESC[QUERY_OPERATOR.CONTAINS_CS]
}

export const QUERY_OPERATOR_OTHER_DESC = {
  ...QUERY_OPERATOR_DESC,
  [QUERY_OPERATOR.LIKE]: QUERY_OPERATOR_DESC[QUERY_OPERATOR.CONTAINS]
}

const mapping = {
  [QUERY_OPERATOR.EQ]: 'equal',
  [QUERY_OPERATOR.NE]: 'not_equal',
  [QUERY_OPERATOR.IN]: 'in',
  [QUERY_OPERATOR.NIN]: 'not_in',
  [QUERY_OPERATOR.LT]: 'less',
  [QUERY_OPERATOR.LTE]: 'less_or_equal',
  [QUERY_OPERATOR.GT]: 'greater',
  [QUERY_OPERATOR.GTE]: 'greater_or_equal',
  [QUERY_OPERATOR.RANGE]: 'between',
  [QUERY_OPERATOR.NRANGE]: 'not_between',
  [QUERY_OPERATOR.LIKE]: 'contains',
  [QUERY_OPERATOR.CONTAINS]: 'contains',
  [QUERY_OPERATOR.CONTAINS_CS]: 'contains_s'
}

// 主机预览接口特殊处理
export const TRANSFORM_SPECIAL_HANDLE_OPERATOR = {
  [QUERY_OPERATOR_SYMBOL[QUERY_OPERATOR.CONTAINS]]: QUERY_OPERATOR.LIKE,
  [QUERY_OPERATOR.LIKE]: mapping[QUERY_OPERATOR.CONTAINS_CS]
}
// 动态分组集群新建/编辑接口特殊处理
export const TRANSFORM_DYNAMIC_SET_OPERATOR = {
  [QUERY_OPERATOR.LIKE]: mapping[QUERY_OPERATOR.CONTAINS]
}

// 动态分组详情接口回显操作符处理(查询对象为集群，对象为主机的主机字段属性)
export const OPERATOR_ECHO = {
  [mapping[QUERY_OPERATOR.CONTAINS_CS]]: QUERY_OPERATOR.CONTAINS_CS,
  [mapping[QUERY_OPERATOR.CONTAINS]]: QUERY_OPERATOR.CONTAINS
}
// 动态分组详情接口回显操作符处理(查询对象为主机的非主机字段属性)
export const OPERATOR_SPECIAL_ECHO = {
  [QUERY_OPERATOR.LIKE]: QUERY_OPERATOR.CONTAINS,
  [mapping[QUERY_OPERATOR.CONTAINS_CS]]: QUERY_OPERATOR.LIKE
}

export default operator => mapping[operator]
