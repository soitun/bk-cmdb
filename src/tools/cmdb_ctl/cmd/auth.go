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

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"configcenter/src/ac"
	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewAuthCommand())
}

type authConf struct {
	resource     string
	resourceFile string
	env          string
	logv         int32
}

// NewAuthCommand TODO
func NewAuthCommand() *cobra.Command {
	conf := new(authConf)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "auth operations",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	subCmds := make([]*cobra.Command, 0)

	userName := new(string)
	supplierAccount := new(string)
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "check if user has the authority to operate resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuthCheckCmd(conf, *userName, *supplierAccount)
		},
	}
	checkCmd.Flags().StringVar(userName, "user", "", "the name of the user")
	checkCmd.Flags().StringVar(supplierAccount, "supplier-account", "0", "the supplier id that this user belongs to")
	subCmds = append(subCmds, checkCmd)

	for _, subCmd := range subCmds {
		cmd.AddCommand(subCmd)
	}
	conf.addFlags(cmd)

	return cmd
}

func (c *authConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&c.resource, "resource", "r", "", "the resource for authorize")
	cmd.PersistentFlags().StringVarP(&c.resourceFile, "rsc-file", "f", "", "the resource file path for authorize")
	cmd.PersistentFlags().StringVarP(&c.env, "environment", "e", "", "the environment for service discovery")
	cmd.PersistentFlags().Int32VarP(&c.logv, "logV", "v", 0,
		"the log level of request, default request body log level is 4")
}

type authService struct {
	authorizer ac.AuthorizeInterface
	resource   []meta.ResourceAttribute
}

func newAuthService(c *authConf) (*authService, error) {
	blog.SetV(c.logv)
	if c.resource == "" && c.resourceFile == "" {
		return nil, errors.New("resource must be set via resource flag or resource file specified by rsc-file flag")
	}

	client := zk.NewZkClient(config.Conf.ZkAddr, 40*time.Second, &config.Conf.ZkTLS)
	if err := client.Start(); err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", config.Conf.ZkAddr, err)
	}
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", config.Conf.ZkAddr, err)
	}
	serviceDiscovery, err := discovery.NewServiceDiscovery(client, c.env)
	if err != nil {
		return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", config.Conf.ZkAddr, err)
	}
	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}
	clientSet, err := apimachinery.NewApiMachinery(apiMachineryConfig, serviceDiscovery)
	if err != nil {
		return nil, fmt.Errorf("new api machinery failed, err: %v", err)
	}
	service := &authService{
		authorizer: iam.NewAuthorizer(clientSet),
	}

	if c.resource != "" {
		err = json.Unmarshal([]byte(c.resource), &service.resource)
		if err != nil {
			return nil, err
		}
	} else {
		resourceFile, err := os.Open(c.resourceFile)
		if err != nil {
			return nil, fmt.Errorf("fail to open file(%s), err(%s)", c.resourceFile, err.Error())
		}
		defer resourceFile.Close()
		resource, err := ioutil.ReadAll(resourceFile)
		if err != nil {
			blog.Errorf("fail to read data from resource file(%s), err:%s", resourceFile, err.Error())
			return nil, err
		}
		err = json.Unmarshal(resource, &service.resource)
		if err != nil {
			return nil, err
		}
	}
	return service, nil
}

func runAuthCheckCmd(c *authConf, userName string, supplierAccount string) error {
	srv, err := newAuthService(c)
	if err != nil {
		return err
	}
	header := headerutil.BuildHeader("admin", common.BKDefaultOwnerID)

	userInfo := meta.UserInfo{
		UserName:        userName,
		SupplierAccount: supplierAccount,
	}
	decisions, err := srv.authorizer.AuthorizeBatch(context.Background(), header, userInfo, srv.resource...)
	if err != nil {
		return err
	}
	for _, decision := range decisions {
		if !decision.Authorized {
			_, _ = fmt.Fprintln(os.Stdout, WithGreenColor("Unauthorized"))
			return nil
		}
	}
	_, _ = fmt.Fprintln(os.Stdout, WithGreenColor("Authorized"))
	return nil
}
