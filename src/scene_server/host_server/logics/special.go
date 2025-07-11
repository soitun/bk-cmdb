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

package logics

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

/*
文件描述： 该文件文件中的接口是专用接口，对应解决某个问题API。 不对UI,第三方应用开放。

*/

// special 用来做隔离， 不让logic 正常逻辑调用
type special struct {
	kit *rest.Kit
	lgc *Logics
}

// SpecialHandle 用来做隔离， 不让logic 正常逻辑调用
type SpecialHandle interface {
	BkSystemInstall(ctx context.Context, appName string, input *metadata.BkSystemInstallRequest) errors.CCError
}

// NewSpecial return handle special logic
func (lgc *Logics) NewSpecial(kit *rest.Kit) SpecialHandle {
	return &special{
		kit: kit,
		lgc: lgc,
	}
}

// BkSystemInstall 蓝鲸组件机器安装agent，主机写入cmdb
// 描述: 1. 不能将主机转移到空闲机和故障机等内置模块
// 2. 不会删除主机已经存在的主机模块， 只会新加主机与模块。 3. 不存在的主机会新加， 规则通过内网IP和 cloud id 判断主机是否存在
// 4. 进程不存在不报错
func (s *special) BkSystemInstall(ctx context.Context, appName string,
	input *metadata.BkSystemInstallRequest) errors.CCError {

	if input.HostInfo == nil {
		input.HostInfo = make(map[string]interface{}, 0)
	}
	// handle host info
	input.HostInfo[common.BKHostInnerIPField] = input.InnerIP
	input.HostInfo[common.BKCloudIDField] = input.CloudID

	// bkSystemParameterConv  将名字转为cc id， 主机返回值的方式和hostID. hostID=0,表示主机不存
	appID, moduleIDArr, hostID, err := s.bkSystemParameterConv(ctx, appName, input)
	if err != nil {
		blog.ErrorJSON("BkSystemInstall convert name to cc id error. err:%s, input:%s, rid:%s", err, input, s.kit.Rid)
		return err
	}
	// TODO auth logic
	// 这里后续需要加鉴权方式
	// 先鉴权后操作

	if hostID == 0 {
		// host not found
		hostID, err = s.bkSystemInstallAddHostInstance(ctx, input)
		if err != nil {
			blog.Errorf("BkSystemInstall IsHostExistInApp error. err:%s, parameters:%s, rid:%s", err.Error(), input,
				s.kit.Rid)
			return err
		}
	} else {
		// check host belong app
		// source host belong app
		ok, err := s.lgc.IsHostExistInApp(s.kit, appID, hostID)
		if err != nil {
			blog.Errorf("BkSystemInstall IsHostExistInApp error. err:%s, params:{appID:%d, hostID:%d}, rid:%s",
				err.Error(), hostID, s.kit.Rid)
			return err
		}
		if !ok {
			blog.Errorf("BkSystemInstall Host does not belong to the current application; error, params:{appID:%d, "+
				"hostID:%d}, rid:%s", appID, hostID, s.kit.Rid)
			return s.kit.CCError.CCErrorf(common.CCErrHostNotINAPP, hostID)
		}

		updateInput := &metadata.UpdateOption{
			Data:      input.HostInfo,
			Condition: mapstr.MapStr{common.BKHostIDField: hostID},
		}
		_, httpDoErr := s.lgc.CoreAPI.CoreService().Instance().UpdateInstance(ctx, s.kit.Header,
			common.BKInnerObjIDHost, updateInput)
		if httpDoErr != nil {
			blog.ErrorJSON("BkSystemInstall update host instance http do error.  err:%s, input:%s, "+
				"update parameter:%s, rid:%s", httpDoErr, input, updateInput, s.kit.Rid)
			return httpDoErr
		}
	}

	err = s.bkSystemInstallModule(ctx, appID, hostID, moduleIDArr)
	if err != nil {
		blog.ErrorJSON("BkSystemInstallBkSystemInstall bkSystemInstallModule error. err:%s, parameters:%s, rid:%s",
			err.Error(), input, s.kit.Rid)
		return err
	}

	// 进程不存在不报错
	err = s.bkSystemInstallProc(ctx, appID, moduleIDArr, hostID, input.ProcInfo)
	if err != nil {
		blog.Errorf("BkSystemInstallBkSystemInstall bkSystemInstallProc error. err:%s, parameters:%s, rid:%s",
			err.Error(), input, s.kit.Rid)
		return err
	}
	return nil
}

// bkSystemParameterConv  将名字转为cc id， 主机返回值的方式和hostID. hostID=0,表示主机不存
func (s *special) bkSystemParameterConv(ctx context.Context, appName string, input *metadata.BkSystemInstallRequest) (appID int64, moduleIDArr []int64, hostID int64, err errors.CCError) {
	bkAppCond := []metadata.ConditionItem{
		metadata.ConditionItem{
			Field:    common.BKAppNameField,
			Operator: common.BKDBEQ,
			Value:    appName,
		},
	}

	var appIDArr []int64
	appIDArr, err = s.lgc.GetAppIDByCond(s.kit, metadata.ConditionWithTime{Condition: bkAppCond})
	if err != nil {
		blog.ErrorJSON("bkSystemParameterConv get blueking app error. err:%s, cond:%s, rid:%s", err.Error(), bkAppCond, s.kit.Rid)
		return
	}

	if len(appIDArr) == 0 {
		blog.ErrorJSON("bkSystemParameterConv blueking app not found. cond:%s, rid:%s", bkAppCond, s.kit.Rid)
		err = s.kit.CCError.CCErrorf(common.CCErrCommBizNotFoundError, appName)
		return
	}
	appID = appIDArr[0]

	moduleIDArr, err = s.bkSystemGetInstallModuleID(ctx, appID, input.SetName, input.ModuleName)
	if err != nil {
		blog.ErrorJSON("bkSystemParameterConv bkSystemGetInstallModuleID error. err:%s, input:%s, rid:%s", err.Error(), input, s.kit.Rid)
		return
	}
	if len(moduleIDArr) == 0 {
		blog.ErrorJSON("bkSystemParameterConv bkSystemGetInstallModuleID not found module id. input:%s, rid:%s", input, s.kit.Rid)
		err = s.kit.CCError.Errorf(common.CCErrCommTopoModuleNotFoundError, input.SetName+"->"+input.ModuleName)
		return
	}

	isExist, err := s.lgc.IsPlatAllExist(s.kit, []int64{input.CloudID})
	if nil != err {
		blog.ErrorJSON("bkSystemParameterConv get cloud  error. err:%s, cond:%s, rid:%s", err.Error(), input, s.kit.Rid)
		return
	}
	if !isExist {
		err = s.kit.CCError.Error(common.CCErrTopoCloudNotFound)
		return
	}

	_, hostID, err = s.lgc.IPCloudToHost(s.kit, input.InnerIP, input.CloudID)
	if err != nil {
		blog.InfoJSON("bkSystemParameterConv IPCloudToHost error. err:%s, input:%s, rid:%s", err, input, s.kit.Rid)
		return
	}

	return

}

// bkSystemInstallAddHostInstance only add host instance. not add host and module relation
func (s *special) bkSystemInstallAddHostInstance(ctx context.Context, input *metadata.BkSystemInstallRequest) (int64,
	errors.CCError) {

	resp, httpDoErr := s.lgc.CoreAPI.CoreService().Instance().
		CreateInstance(ctx, s.kit.Header, common.BKInnerObjIDHost, &metadata.CreateModelInstance{
			Data: input.HostInfo,
		})
	if httpDoErr != nil {
		blog.ErrorJSON("BkSystemInstall create host instance http do error.  err:%s, data:%s, rid:%s", httpDoErr,
			input.HostInfo, s.kit.Rid)
		return 0, s.kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	return int64(resp.Created.ID), nil

}

// bkSystemGetInstallModuleID only add host instance. not add host and module relation
func (s *special) bkSystemGetInstallModuleID(ctx context.Context, appID int64, setName, moduleName string) ([]int64, errors.CCError) {
	bkSetCond := []metadata.ConditionItem{
		metadata.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    appID,
		},
		metadata.ConditionItem{
			Field:    common.BKSetNameField,
			Operator: common.BKDBEQ,
			Value:    setName,
		},
	}

	setIDArr, err := s.lgc.GetSetIDByCond(s.kit, metadata.ConditionWithTime{Condition: bkSetCond})
	if err != nil {
		blog.ErrorJSON("bkSystemGetInstallModuleID GetSetIDByCond error. err:%s, cond:%s, rid:%s", err.Error(), bkSetCond, s.kit.Rid)
		return nil, err
	}
	if len(setIDArr) == 0 {
		blog.Warnf("bkSystemGetInstallModuleID GetSetIDByCond not found set. cond:%v, rid:%s", bkSetCond, s.kit.Rid)
		return nil, nil
	}

	bkModuleCond := []metadata.ConditionItem{
		metadata.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    appID,
		},
		metadata.ConditionItem{
			Field:    common.BKSetIDField,
			Operator: common.BKDBIN,
			Value:    setIDArr,
		},
		metadata.ConditionItem{
			Field:    common.BKModuleNameField,
			Operator: common.BKDBEQ,
			Value:    moduleName,
		},
	}
	moduleIDArr, err := s.lgc.GetModuleIDByCond(s.kit, metadata.ConditionWithTime{Condition: bkModuleCond})
	if err != nil {
		blog.ErrorJSON("bkSystemGetInstallModuleID GetModuleIDByCond error. err:%s, cond:%s, rid:%s", err.Error(), bkModuleCond, s.kit.Rid)
		return nil, err
	}
	if len(moduleIDArr) == 0 {
		blog.Warnf("bkSystemGetInstallModuleID GetModuleIDByCond not found set. cond:%v, rid:%s", bkModuleCond, s.kit.Rid)
	}
	return moduleIDArr, nil
}

// bkSystemInstallModule TODO
// bksystemInstallAddHostInstance only all host and module relation
func (s *special) bkSystemInstallModule(ctx context.Context, appID, hostID int64, moduleIDArr []int64) errors.CCError {
	input := &metadata.HostsModuleRelation{
		ApplicationID: appID,
		HostID:        []int64{hostID},
		ModuleID:      moduleIDArr,
		IsIncrement:   true,
	}

	resp, err := s.lgc.CoreAPI.CoreService().Host().TransferToNormalModule(ctx, s.kit.Header, input)
	if err != nil {
		blog.Errorf("transfer host failed, err: %v, resp: %#v, input: %#v, rid:%s", err, resp, input, s.kit.Rid)
		return err
	}

	return nil
}

// bkSystemInstallProc  change host process info. process not found,
// process does not exist without error. no performance issues are consider,
func (s *special) bkSystemInstallProc(ctx context.Context, appID int64, moduleIDArr []int64, hostID int64,
	procInfoMap map[string]map[string]interface{}) errors.CCError {

	for key := range procInfoMap {
		delete(procInfoMap[key], common.BKFuncName)
		delete(procInfoMap[key], common.BKProcessNameField)
		delete(procInfoMap[key], common.BKProcessIDField)
		delete(procInfoMap[key], common.CreateTimeField)
	}
	for _, moduleID := range moduleIDArr {
		searchSrvInstRelationCond := &metadata.ListServiceInstanceDetailOption{
			BusinessID: appID,
			HostID:     hostID,
			ModuleID:   moduleID,
		}

		srvInstInfo, err := s.lgc.CoreAPI.CoreService().Process().ListServiceInstanceDetail(ctx, s.kit.Header,
			searchSrvInstRelationCond)
		if err != nil {
			blog.ErrorJSON("bkSystemInstallProc ListServiceInstance  http  error.  err:%s, data:%s, rid:%s", err,
				appID, searchSrvInstRelationCond, s.kit.Rid)
			return err
		}
		var procIDArr []int64
		for _, srvInst := range srvInstInfo.Info {
			for _, procInfo := range srvInst.ProcessInstances {
				procIDArr = append(procIDArr, procInfo.Process.ProcessID)
			}

		}
		for procName, info := range procInfoMap {
			updateCond := condition.CreateCondition()
			updateCond.Field(common.BKProcessIDField).In(procIDArr)
			updateCond.Field(common.BKFuncName).Eq(procName)

			procUpdateOpt := &metadata.UpdateOption{
				Data:      info,
				Condition: updateCond.ToMapStr(),
			}
			_, httpDoErr := s.lgc.CoreAPI.CoreService().Instance().UpdateInstance(ctx, s.kit.Header,
				common.BKInnerObjIDProc, procUpdateOpt)
			if httpDoErr != nil {
				blog.ErrorJSON("bkSystemInstallProc UpdateInstance  http do error.  err:%s, data:%s, rid:%s",
					httpDoErr, appID, searchSrvInstRelationCond, s.kit.Rid)
				return httpDoErr
			}
		}

	}

	return nil
}
