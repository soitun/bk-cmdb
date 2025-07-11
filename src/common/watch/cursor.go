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

package watch

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/stream/types"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NoEventCursor TODO
// this cursor means there is no event occurs.
// we send this cursor to our the watcher and if we
// received a NoEventCursor, then we need to fetch event
// from the head cursor
var NoEventCursor string

func init() {
	no := Cursor{
		Type:        NoEvent,
		ClusterTime: types.TimeStamp{Sec: 1, Nano: 1},
		Oid:         "5ea6d3f394c1f5d986e9bd86",
		Oper:        types.OperType("noEvent"),
	}
	cursor, err := no.Encode()
	if err != nil {
		panic("initial NoEventCursor failed")
	}
	// cursor should be:
	// MQ0xDTVlYTZkM2YzOTRjMWY1ZDk4NmU5YmQ4Ng1ub0V2ZW50DTENMQ==
	NoEventCursor = cursor

	for cursorType, intVal := range cursorTypeIntMap {
		intCursorTypeMap[intVal] = cursorType
	}
}

var (
	cursorTypeIntMap = map[CursorType]int{
		NoEvent:                 1,
		Host:                    2,
		ModuleHostRelation:      3,
		Biz:                     4,
		Set:                     5,
		Module:                  6,
		ObjectBase:              8,
		Process:                 9,
		ProcessInstanceRelation: 10,
		HostIdentifier:          11,
		MainlineInstance:        12,
		InstAsst:                13,
		BizSet:                  14,
		BizSetRelation:          15,
		Plat:                    16,
		KubeCluster:             17,
		KubeNode:                18,
		KubeNamespace:           19,
		KubeWorkload:            20,
		KubePod:                 21,
		Project:                 22,
	}

	intCursorTypeMap = make(map[int]CursorType)
)

// CursorType TODO
type CursorType string

const (
	// NoEvent TODO
	NoEvent CursorType = "no_event"
	// UnknownType TODO
	UnknownType CursorType = "unknown"
	// Host TODO
	Host CursorType = "host"
	// ModuleHostRelation TODO
	ModuleHostRelation CursorType = "host_relation"
	// Biz TODO
	Biz CursorType = "biz"
	// Set TODO
	Set CursorType = "set"
	// Module TODO
	Module CursorType = "module"
	// ObjectBase TODO
	ObjectBase CursorType = "object_instance"
	// Process TODO
	Process CursorType = "process"
	// ProcessInstanceRelation TODO
	ProcessInstanceRelation CursorType = "process_instance_relation"
	// BizSet TODO
	BizSet CursorType = "biz_set"
	// HostIdentifier TODO
	// a mixed event type, which contains host, host relation, process events etc.
	// and finally converted to host identifier event.
	HostIdentifier CursorType = "host_identifier"
	// MainlineInstance specified for mainline instance event watch, filtered from object instance events
	MainlineInstance CursorType = "mainline_instance"
	// InstAsst specified for instance association event watch
	InstAsst CursorType = "inst_asst"
	// BizSetRelation a mixed event type containing biz set & biz events, which are converted to their relation events
	BizSetRelation CursorType = "biz_set_relation"
	// Plat cloud area event cursor type
	Plat CursorType = "plat"
	// Project project event cursor type
	Project CursorType = "project"
	// kube related cursor types
	// KubeCluster cursor type
	KubeCluster CursorType = "kube_cluster"
	// KubeNode cursor type
	KubeNode CursorType = "kube_node"
	// KubeNamespace cursor type
	KubeNamespace CursorType = "kube_namespace"
	// KubeWorkload cursor type, including all workloads(e.g. deployment) with their type specified in sub-resource
	KubeWorkload CursorType = "kube_workload"
	// KubePod cursor type, its event detail is pod info with containers in it
	KubePod CursorType = "kube_pod"
)

// ToInt TODO
func (ct CursorType) ToInt() int {
	intVal, exists := cursorTypeIntMap[ct]
	if !exists {
		return -1
	}

	return intVal
}

// ParseInt TODO
func (ct *CursorType) ParseInt(typ int) {
	cursorType, exists := intCursorTypeMap[typ]
	if !exists {
		*ct = UnknownType
		return
	}

	*ct = cursorType
}

// ListCursorTypes returns all support CursorTypes.
func ListCursorTypes() []CursorType {
	return []CursorType{Host, ModuleHostRelation, Biz, Set, Module, ObjectBase, Process, ProcessInstanceRelation,
		HostIdentifier, MainlineInstance, InstAsst, BizSet, BizSetRelation, Plat, KubeCluster, KubeNode, KubeNamespace,
		KubeWorkload, KubePod, Project}
}

// Cursor is a self-defined token which is corresponding to the mongodb's resume token.
// cursor has a unique and 1:1 relationship with mongodb's resume token.
type Cursor struct {
	Type        CursorType
	ClusterTime types.TimeStamp
	// a random hex string to avoid the caller to generated a self-defined cursor.
	Oid  string
	Oper types.OperType
	// UniqKey is an optional key which is used to ensure that a cursor is unique among a same resources(
	// as is Type field).
	// This key is used when the upper fields can not describe a unique cursor. such as the common object instance
	// event, all these instance event all have a same cursor type, as is object instance. but the instance id is
	// unique which can be used as a unique key, and is REENTRANT when a retry operation happens which is
	// **VERY IMPORTANT**.
	UniqKey string
}

const cursorVersion = "1"

// Encode TODO
func (c Cursor) Encode() (string, error) {
	if c.Type == "" {
		return "", errors.New("unsupported type")
	}

	if c.ClusterTime.Sec == 0 {
		return "", errors.New("invalid cluster time sec")
	}

	if c.Oid == "" {
		return "", errors.New("invalid oid")
	}

	if c.Oper == "" {
		return "", errors.New("unsupported operation type")
	}

	sec := strconv.FormatUint(uint64(c.ClusterTime.Sec), 10)
	nano := strconv.FormatUint(uint64(c.ClusterTime.Nano), 10)
	pool := bytes.Buffer{}
	// version field.
	pool.WriteString(cursorVersion)
	pool.WriteByte('\r')

	// type filed.
	if c.Type.ToInt() < 0 {
		return "", errors.New("unsupported cursor type")
	}

	pool.WriteString(strconv.Itoa(c.Type.ToInt()))
	pool.WriteByte('\r')

	// oid field.
	pool.WriteString(c.Oid)
	pool.WriteByte('\r')

	// operation type field
	pool.WriteString(string(c.Oper))
	pool.WriteByte('\r')

	// cluster time sec field.
	pool.WriteString(sec)
	pool.WriteByte('\r')

	// cluster time nano field
	pool.WriteString(nano)
	pool.WriteByte('\r')

	// unique key field
	pool.WriteString(c.UniqKey)

	return base64.StdEncoding.EncodeToString(pool.Bytes()), nil
}

// Decode TODO
func (c *Cursor) Decode(cur string) error {
	byt, err := base64.StdEncoding.DecodeString(cur)
	if err != nil {
		return fmt.Errorf("decode cursor, but base64 decode failed, err: %v", err)
	}

	elements := make([]string, 0)
	pool := bytes.NewBuffer(byt)

	ele := make([]byte, 0)
	for {
		b, err := pool.ReadByte()
		if err != nil {
			if err != io.EOF {
				return err
			}
			// to the end
			elements = append(elements, string(ele))
			break
		}
		if b == '\r' {
			elements = append(elements, string(ele))
			ele = ele[:0]
		} else {
			ele = append(ele, b)
		}
	}

	// at least have 6 fields, UniqKey is an optional field.
	if len(elements) < 6 {
		return errors.New("invalid cursor string")
	}

	if elements[0] != cursorVersion {
		return fmt.Errorf("decode cursor, but got invalid cursor version: %s", elements[0])
	}

	typ, err := strconv.Atoi(elements[1])
	if err != nil {
		return fmt.Errorf("got invalid type: %s", elements[1])
	}
	cursorType := CursorType("")
	cursorType.ParseInt(typ)
	c.Type = cursorType

	switch cursorType {
	case InstAsst:
		// instance association events use its identity id which is formatted to a string type as its event's oid.
		// so its oid should be validate with ParseInt function.
		_, err = strconv.ParseInt(elements[2], 10, 64)
		if err != nil {
			return fmt.Errorf("got invalid oid: %s, should be a string formatted from int, err: %v", elements[2], err)
		}
	default:
		_, err = primitive.ObjectIDFromHex(elements[2])
		if err != nil {
			return fmt.Errorf("got invalid oid: %s, should be a hex string, err: %v", elements[2], err)
		}
	}
	c.Oid = elements[2]

	if elements[3] == "" {
		return fmt.Errorf("decode cursor, but got empty operation type")
	}
	c.Oper = types.OperType(elements[3])

	sec, err := strconv.ParseUint(elements[4], 10, 64)
	if err != nil {
		return fmt.Errorf("got invalid sec field %s, err: %v", elements[4], err)
	}
	c.ClusterTime.Sec = uint32(sec)

	nano, err := strconv.ParseUint(elements[5], 10, 64)
	if err != nil {
		return fmt.Errorf("got invalid nano field %s, err: %v", elements[5], err)
	}
	c.ClusterTime.Nano = uint32(nano)

	// cause unique key is an optional key.
	if len(elements) >= 7 {
		c.UniqKey = elements[6]
	}

	return nil
}

var collEventCursorTypeMap = map[string]CursorType{
	common.BKTableNameBaseHost:                Host,
	common.BKTableNameModuleHostConfig:        ModuleHostRelation,
	common.BKTableNameBaseApp:                 Biz,
	common.BKTableNameBaseSet:                 Set,
	common.BKTableNameBaseModule:              Module,
	common.BKTableNameBaseInst:                ObjectBase,
	common.BKTableNameMainlineInstance:        MainlineInstance,
	common.BKTableNameBaseProcess:             Process,
	common.BKTableNameProcessInstanceRelation: ProcessInstanceRelation,
	common.BKTableNameInstAsst:                InstAsst,
	common.BKTableNameBaseBizSet:              BizSet,
	common.BKTableNameBasePlat:                Plat,
	kubetypes.BKTableNameBaseCluster:          KubeCluster,
	kubetypes.BKTableNameBaseNode:             KubeNode,
	kubetypes.BKTableNameBaseNamespace:        KubeNamespace,
	kubetypes.BKTableNameBaseWorkload:         KubeWorkload,
	kubetypes.BKTableNameBasePod:              KubePod,
	common.BKTableNameBaseProject:             Project,
}

// GetEventCursor get event cursor.
func GetEventCursor(coll string, e *types.Event, instID int64) (string, error) {
	curType, exists := collEventCursorTypeMap[coll]
	if !exists {
		blog.Errorf("unsupported cursor type collection: %s, oid: %s", e.ID())
		return "", fmt.Errorf("unsupported cursor type collection: %s", coll)
	}

	hCursor := &Cursor{
		Type:        curType,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		Oper:        e.OperationType,
	}

	switch curType {
	case ObjectBase, MainlineInstance, InstAsst:
		if instID <= 0 {
			return "", errors.New("invalid instance id")
		}

		// add unique key for common object instance.
		hCursor.UniqKey = strconv.FormatInt(instID, 10)
	case KubeWorkload:
		if instID <= 0 {
			return "", errors.New("invalid kube workload id")
		}

		// add unique key for kube workload, composed by workload type and id.
		hCursor.UniqKey = fmt.Sprintf("%s:%d", gjson.GetBytes(e.DocBytes, kubetypes.KindField).String(), instID)
	}

	hCursorEncode, err := hCursor.Encode()
	if err != nil {
		blog.Errorf("encode node cursor failed, err: %v, oid: %s", err, e.Oid)
		return "", err
	}

	return hCursorEncode, nil
}
