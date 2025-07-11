/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package operator

// Calculator is used to calculate the logic value of a list/array of
// Operator with different of "Calculator Type" instance.
type Calculator interface {
	// Name of the calculator
	Name() string

	// Result TODO
	// the calculated resulted with multiple Policy.
	Result(p []*Policy) (bool, error)
}

// And TODO
const And = OperType("AND")

// AndOper TODO
type AndOper OperType

// Name TODO
func (a *AndOper) Name() string {
	return "AND"
}

// Result TODO
func (a *AndOper) Result(p []*Policy) (bool, error) {
	return true, nil
}

// Or TODO
const Or = OperType("OR")

// OrOper TODO
type OrOper OperType

// Name TODO
func (o *OrOper) Name() string {
	return "OR"
}
