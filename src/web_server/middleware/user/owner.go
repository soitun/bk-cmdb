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

package user

import (
	"context"
	"net/http"
	"time"

	"configcenter/src/apimachinery/apiserver"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	validator "configcenter/src/source_controller/coreservice/core/instances"
	"configcenter/src/storage/dal/redis"
)

// OwnerManager TODO
type OwnerManager struct {
	Engine   *backbone.Engine
	CacheCli redis.Client
	OwnerID  string
	UserName string
	header   http.Header
	ApiCli   apiserver.ApiServerClientInterface
}

// NewOwnerManager TODO
func NewOwnerManager(userName, ownerID, language string) *OwnerManager {
	ownerManager := new(OwnerManager)
	ownerManager.UserName = userName
	ownerManager.OwnerID = ownerID

	header := headerutil.BuildHeader(userName, ownerID)
	httpheader.SetLanguage(header, language)
	ownerManager.header = header

	return ownerManager
}

// SetHttpHeader TODO
func (m *OwnerManager) SetHttpHeader(key, val string) {
	m.header.Set(key, val)
}

// InitOwner TODO
func (m *OwnerManager) InitOwner() (*metadata.IamPermission, errors.CCErrorCoder) {
	rid := httpheader.GetRid(m.header)
	blog.V(5).Infof("init owner %s, rid: %s", m.OwnerID, rid)
	ccErr := m.Engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(m.header))

	exist, err, permissions := m.defaultAppIsExist()
	if err != nil {
		return permissions, err
	}
	if !exist {
		redisCli := m.CacheCli
		for {
			ok, err := redisCli.SetNX(context.Background(), common.BKCacheKeyV3Prefix+"owner_init_lock:"+m.OwnerID,
				m.OwnerID, 60*time.Second).Result()
			if nil != err {
				blog.Errorf("owner_init_lock error %s, rid: %s", err.Error(), rid)
				return nil, ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
			}
			if ok {
				break
			}
			time.Sleep(time.Second)
		}
		defer func() {
			if err := redisCli.Del(context.Background(),
				common.BKCacheKeyV3Prefix+"owner_init_lock:"+m.OwnerID).Err(); err != nil {
				blog.Errorf("owner_init_lock error %s, rid: %s", err.Error(), rid)
			}
		}()
		exist, err, permissions = m.defaultAppIsExist()
		if err != nil {
			return permissions, err
		}
		if !exist {
			err, permissions = m.addDefaultApp()
			if nil != err {
				return permissions, err
			}
		}
	}
	return nil, nil
}

func (m *OwnerManager) addDefaultApp() (errors.CCErrorCoder, *metadata.IamPermission) {
	rid := httpheader.GetRid(m.header)
	ccErr := m.Engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(m.header))

	blog.V(5).Infof("addDefaultApp %s, rid: %s", m.OwnerID, rid)
	params, err, permissions := m.getObjectFields(common.BKInnerObjIDApp)
	if err != nil {
		return err, permissions
	}
	params[common.BKAppNameField] = common.DefaultAppName
	params[common.BKMaintainersField] = "admin"
	params[common.BKProductPMField] = "admin"
	params[common.BKTimeZoneField] = "Asia/Shanghai"
	params[common.BKLanguageField] = "1" // 中文
	params[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal

	result, httpDoErr := m.ApiCli.AddDefaultApp(context.Background(), m.header, m.OwnerID, params)
	if httpDoErr != nil {
		blog.ErrorJSON("addDefaultApp searchDefaultApp http do error. err:%s, rid:%s", httpDoErr.Error(), rid)
		return ccErr.CCError(common.CCErrCommHTTPDoRequestFailed), nil
	}

	if err := result.CCError(); err != nil {
		blog.ErrorJSON("addDefaultApp searchDefaultApp http replay error. err:%s, rid:%s", result, rid)
		return err, result.Permissions
	}
	return nil, nil
}

func (m *OwnerManager) defaultAppIsExist() (bool, errors.CCErrorCoder, *metadata.IamPermission) {
	ccErr := m.Engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(m.header))
	rid := httpheader.GetRid(m.header)
	result, httpDoErr := m.ApiCli.SearchDefaultApp(context.Background(), m.header, m.OwnerID)
	if httpDoErr != nil {
		blog.ErrorJSON("defaultAppIsExist searchDefaultApp http do error. err:%s, rid:%s", httpDoErr.Error(), rid)
		return false, ccErr.CCError(common.CCErrCommHTTPDoRequestFailed), nil
	}

	if err := result.CCError(); err != nil {
		blog.ErrorJSON("defaultAppIsExist searchDefaultApp http replay error. err:%s, rid:%s", result, rid)
		return false, err, result.Permissions
	}

	return 0 < result.Data.Count, nil, nil
}

func (m *OwnerManager) getObjectFields(objID string) (map[string]interface{}, errors.CCErrorCoder,
	*metadata.IamPermission) {

	ccErr := m.Engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(m.header))
	rid := httpheader.GetRid(m.header)

	filter := mapstr.MapStr{
		common.BKObjIDField: objID,
		"page": common.KvMap{
			"skip":  0,
			"limit": common.BKNoLimit,
		},
	}
	result, httpDoErr := m.ApiCli.GetObjectAttr(context.Background(), m.header, filter)
	if httpDoErr != nil {
		blog.Errorf("get object attribute failed, err: %v, cond: %+v, rid: %s", httpDoErr, filter, rid)
		return nil, ccErr.CCError(common.CCErrCommHTTPDoRequestFailed), nil
	}

	if err := result.CCError(); err != nil {
		blog.Errorf("get object attribute failed, err: %v, cond: %+v, rid: %s", err, filter, rid)
		return nil, err, result.Permissions
	}

	fields := result.Data

	ret := map[string]interface{}{}
	validator.FillLostFieldValue(context.Background(), ret, fields)
	return ret, nil, nil
}
