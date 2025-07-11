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

package configcenter

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	ccerr "configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/ssl"
	"configcenter/src/storage/dal/kafka"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/viper"
)

// configuration parser representing five files,the configuration can be taken out according to the parser.
var redisParser *viperParser
var mongodbParser *viperParser
var commonParser *viperParser
var extraParser *viperParser
var localParser *viperParser

var confLock sync.RWMutex

func checkDir(path string) error {
	info, err := os.Stat(path)
	if os.ErrNotExist == err {
		return fmt.Errorf("directory %s not exists", path)
	}
	if err != nil {
		return fmt.Errorf("stat directory %s faile, %s", path, err.Error())
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not directory", path)
	}

	return nil
}

func loadErrorAndLanguage(errorres string, languageres string, handler *CCHandler) error {
	if err := checkDir(errorres); err != nil {
		return err
	}
	errcode, err := ccerr.LoadErrorResourceFromDir(errorres)
	if err != nil {
		return fmt.Errorf("load error resource error: %s", err)
	}
	handler.OnErrorUpdate(nil, errcode)

	if err := checkDir(languageres); err != nil {
		return err
	}
	languagepack, err := language.LoadLanguageResourceFromDir(languageres)
	if err != nil {
		return fmt.Errorf("load language resource error: %s", err)
	}
	handler.OnLanguageUpdate(nil, languagepack)
	return nil
}

// GetLocalConf get local config
func GetLocalConf(confPath string, handler *CCHandler) error {
	if err := SetLocalFile(confPath); err != nil {
		return fmt.Errorf("parse config file: %s failed, err: %v", confPath, err)
	}

	// load local error and language
	errorres, _ := String("errors.res")
	languageres, _ := String("language.res")
	if err := loadErrorAndLanguage(errorres, languageres, handler); err != nil {
		return err
	}

	if handler.OnProcessUpdate != nil {
		handler.OnProcessUpdate(ProcessConfig{}, ProcessConfig{})
	}

	return nil
}

// SetRedisFromByte TODO
func SetRedisFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	if redisParser != nil {
		err := redisParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from redis")
			return err
		}
		return nil
	}
	redisParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from redis")
		return err
	}
	return nil
}

// SetRedisFromFile TODO
func SetRedisFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	redisParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from redis")
		return err
	}
	return nil
}

// SetMongodbFromByte TODO
func SetMongodbFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	if mongodbParser != nil {
		err = mongodbParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from mongodb")
			return err
		}
		return nil
	}
	mongodbParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from mongodb")
		return err
	}
	return nil
}

// SetMongodbFromFile TODO
func SetMongodbFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	mongodbParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from mongodb")
		return err
	}
	return nil
}

// SetCommonFromByte TODO
func SetCommonFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// if it is not nil, do not create a new parser, but add the new configuration information to viper
	if commonParser != nil {
		err = commonParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from common")
			return err
		}
		return nil
	}
	commonParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from common")
		return err
	}
	return nil
}

// SetCommonFromFile TODO
func SetCommonFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	commonParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from common")
		return err
	}
	return nil
}

// SetExtraFromByte TODO
func SetExtraFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// if it is not nil, do not create a new parser, but add the new configuration information to viper
	if extraParser != nil {
		err = extraParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from extra")
			return err
		}
		return nil
	}
	extraParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from extra")
		return err
	}
	return nil
}

// SetExtraFromFile TODO
func SetExtraFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	extraParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from extra")
		return err
	}
	return nil
}

// SetLocalFile set localParser from file
func SetLocalFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// /data/migrate.yaml -> /data/migrate
	split := strings.Split(target, ".")
	filePath := split[0]
	localParser, err = newViperParserFromFile(filePath)
	if err != nil {
		blog.Errorf("set local file failed, target: %s, err: %v", target, err)
		return err
	}
	return nil
}

// Redis return redis configuration information according to the prefix.
func Redis(prefix string) (redis.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	var parser *viperParser
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		parser = getRedisParser()
		if parser != nil {
			break
		}
		blog.Warn("the configuration of redis is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		blog.Errorf("can't find redis configuration")
		return redis.Config{}, errors.New("can't find redis configuration")
	}

	tlsConf, err := NewTLSClientConfigFromConfig(prefix + ".tls")
	if err != nil {
		blog.Errorf("fail to get redis tls configuration")
		return redis.Config{}, err
	}

	return redis.Config{
		Address:          parser.getString(prefix + ".host"),
		Password:         parser.getString(prefix + ".pwd"),
		Database:         parser.getString(prefix + ".database"),
		MasterName:       parser.getString(prefix + ".masterName"),
		SentinelPassword: parser.getString(prefix + ".sentinelPwd"),
		Enable:           parser.getString(prefix + ".enable"),
		MaxOpenConns:     parser.getInt(prefix + ".maxOpenConns"),
		TLSConfig:        &tlsConf,
	}, nil
}

// Mongo return mongo configuration information according to the prefix.
func Mongo(prefix string) (mongo.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	var parser *viperParser
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		parser = getMongodbParser()
		if parser != nil {
			break
		}
		blog.Warn("the configuration of mongo is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		blog.Errorf("can't find mongo configuration")
		return mongo.Config{}, errors.New("can't find mongo configuration")
	}

	tlsClientConfig, err := NewTLSClientConfigFromConfig(prefix + ".tls")
	if err != nil {
		return mongo.Config{}, err
	}

	c := mongo.Config{
		Address:   parser.getString(prefix + ".host"),
		Port:      parser.getString(prefix + ".port"),
		User:      parser.getString(prefix + ".usr"),
		Password:  parser.getString(prefix + ".pwd"),
		Database:  parser.getString(prefix + ".database"),
		Mechanism: parser.getString(prefix + ".mechanism"),
		RsName:    parser.getString(prefix + ".rsName"),
		TLSConf:   &tlsClientConfig,
	}

	if c.RsName == "" {
		blog.Errorf("rsName not set")
	}
	if c.Mechanism == "" {
		c.Mechanism = "SCRAM-SHA-1"
	}

	maxOpenConns := prefix + ".maxOpenConns"
	if !parser.isSet(maxOpenConns) {
		blog.Errorf("can not find config %s, set default value: %d", maxOpenConns, mongo.DefaultMaxOpenConns)
		c.MaxOpenConns = mongo.DefaultMaxOpenConns
	} else {
		c.MaxOpenConns = parser.getUint64(maxOpenConns)
	}

	if c.MaxIdleConns > mongo.MaximumMaxOpenConns {
		blog.Errorf("config %s exceeds maximum value, use maximum value %d", maxOpenConns, mongo.MaximumMaxOpenConns)
		c.MaxIdleConns = mongo.MaximumMaxOpenConns
	}

	maxIdleConns := prefix + ".maxIdleConns"
	if !parser.isSet(maxIdleConns) {
		blog.Errorf("can not find config %s, set default value: %d", maxIdleConns, mongo.MinimumMaxIdleOpenConns)
		c.MaxIdleConns = mongo.MinimumMaxIdleOpenConns
	} else {
		c.MaxIdleConns = parser.getUint64(maxIdleConns)
	}

	if c.MaxIdleConns < mongo.MinimumMaxIdleOpenConns {
		blog.Errorf("config %s less than minimum value, use minimum value %d",
			maxIdleConns, mongo.MinimumMaxIdleOpenConns)
		c.MaxIdleConns = mongo.MinimumMaxIdleOpenConns
	}

	if !parser.isSet(prefix + ".socketTimeoutSeconds") {
		blog.Errorf("can not find mongo.socketTimeoutSeconds config, use default value: %d",
			mongo.DefaultSocketTimeout)
		c.SocketTimeout = mongo.DefaultSocketTimeout
		return c, nil
	}

	c.SocketTimeout = parser.getInt(prefix + ".socketTimeoutSeconds")
	if c.SocketTimeout > mongo.MaximumSocketTimeout {
		blog.Errorf("mongo.socketTimeoutSeconds config %d exceeds maximum value, use maximum value %d",
			c.SocketTimeout, mongo.MaximumSocketTimeout)
		c.SocketTimeout = mongo.MaximumSocketTimeout
	}

	if c.SocketTimeout < mongo.MinimumSocketTimeout {
		blog.Errorf("mongo.socketTimeoutSeconds config %d less than minimum value, use minimum value %d",
			c.SocketTimeout, mongo.MinimumSocketTimeout)
		c.SocketTimeout = mongo.MinimumSocketTimeout
	}

	return c, nil
}

// Kafka return kafka configuration information according to the prefix.
func Kafka(prefix string) (kafka.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	var parser *viperParser
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		parser = getCommonParser()
		if parser != nil {
			break
		}
		blog.Warn("the configuration of common is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		blog.Errorf("can't find kafka configuration")
		return kafka.Config{}, errors.New("can't find kafka configuration")
	}

	return kafka.Config{
		Brokers:   parser.getStringSlice(prefix + ".brokers"),
		GroupID:   parser.getString(prefix + ".groupID"),
		Topic:     parser.getString(prefix + ".topic"),
		Partition: parser.getInt64(prefix + ".partition"),
		User:      parser.getString(prefix + ".user"),
		Password:  parser.getString(prefix + ".password"),
	}, nil
}

// Crypto return crypto configuration information according to the prefix.
func Crypto(prefix string) (*cryptor.Config, error) {
	var parser *viperParser
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		parser = getCommonParser()
		if parser != nil {
			break
		}
		blog.Warn("the configuration of common is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		return nil, errors.New("get common parser failed")
	}

	if !parser.isSet(prefix) {
		return &cryptor.Config{Enabled: false}, nil
	}

	conf := new(cryptor.Config)
	err := parser.unmarshalKey(prefix, conf)
	if err != nil {
		return nil, err
	}

	if conf.Algorithm == "" {
		conf.Algorithm = cryptor.AesGcm
	}

	return conf, nil
}

// String return the string value of the configuration information according to the key.
func String(key string) (string, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return "", err
	}

	return parser.getString(key), nil
}

// Int return the int value of the configuration information according to the key.
func Int(key string) (int, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return 0, err
	}

	if !parser.isConfigIntType(key) {
		return 0, errors.New("config is not int type")
	}
	return parser.getInt(key), nil
}

// Int64 return the int value of the configuration information according to the key.
func Int64(key string) (int64, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return 0, err
	}

	if !parser.isConfigIntType(key) {
		return 0, errors.New("config is not int type")
	}
	return parser.getInt64(key), nil
}

// Bool return the bool value of the configuration information according to the key.
func Bool(key string) (bool, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return false, err
	}

	if !parser.isConfigBoolType(key) {
		return false, errors.New("config is not bool type")
	}
	return parser.getBool(key), nil
}

// StringSlice return the stringSlice value of the configuration information according to the key.
func StringSlice(key string) ([]string, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return nil, err
	}

	return parser.getStringSlice(key), nil
}

// IsExist checks if key exists in all config files
func IsExist(key string) bool {
	confLock.RLock()
	defer confLock.RUnlock()

	// 在所有的配置文件中判断
	if (localParser == nil || !localParser.isSet(key)) && (commonParser == nil || !commonParser.isSet(key)) &&
		(extraParser == nil || !extraParser.isSet(key)) {
		return false
	}
	return true
}

// UnmarshalKey takes a single key and unmarshal it into a Struct.
func UnmarshalKey(key string, val interface{}) error {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return err
	}

	return parser.unmarshalKey(key, val)
}

func getRedisParser() *viperParser {
	if redisParser != nil {
		return redisParser
	}

	return localParser
}

func getMongodbParser() *viperParser {
	if mongodbParser != nil {
		return mongodbParser
	}

	return localParser
}

func getCommonParser() *viperParser {
	if commonParser != nil {
		return commonParser
	}

	return localParser
}

// getKeyValueParser get viper parser for common key value in the order of migrate->common->extra
func getKeyValueParser(key string) (*viperParser, error) {
	if localParser != nil && localParser.isSet(key) {
		return localParser, nil
	}

	if commonParser != nil && commonParser.isSet(key) {
		return commonParser, nil
	}

	if extraParser != nil && extraParser.isSet(key) {
		return extraParser, nil
	}

	if redisParser != nil && redisParser.isSet(key) {
		return redisParser, nil
	}

	if mongodbParser != nil && mongodbParser.isSet(key) {
		return mongodbParser, nil
	}

	return nil, fmt.Errorf("%s key's config not found", key)
}

type viperParser struct {
	parser *viper.Viper
}

func newViperParser(data []byte) (*viperParser, error) {

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return &viperParser{parser: v}, nil
}

func newViperParserFromFile(target string) (*viperParser, error) {
	v := viper.New()
	v.SetConfigName(path.Base(target))
	v.AddConfigPath(path.Dir(target))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	v.WatchConfig()
	return &viperParser{parser: v}, nil
}

func (vp *viperParser) getString(path string) string {
	return vp.parser.GetString(path)
}

func (vp *viperParser) getInt(path string) int {
	return vp.parser.GetInt(path)
}

func (vp *viperParser) getUint64(path string) uint64 {
	return vp.parser.GetUint64(path)
}

func (vp *viperParser) getBool(path string) bool {
	return vp.parser.GetBool(path)
}

func (vp *viperParser) getDuration(path string) time.Duration {
	return vp.parser.GetDuration(path)
}

func (vp *viperParser) isSet(path string) bool {
	return vp.parser.IsSet(path)
}

func (vp *viperParser) getInt64(path string) int64 {
	return vp.parser.GetInt64(path)
}

func (vp *viperParser) getStringSlice(path string) []string {
	return vp.parser.GetStringSlice(path)
}

func (vp *viperParser) isConfigIntType(path string) bool {
	val := vp.parser.GetString(path)
	_, err := strconv.Atoi(val)
	if err != nil {
		return false
	}
	return true
}

func (vp *viperParser) isConfigBoolType(path string) bool {
	val := vp.parser.GetString(path)
	if val != "true" && val != "false" {
		return false
	}
	return true
}

func (vp *viperParser) unmarshalKey(key string, val interface{}) error {
	return vp.parser.UnmarshalKey(key, val)
}

// NewTLSClientConfigFromConfig new config about tls client config
func NewTLSClientConfigFromConfig(prefix string) (ssl.TLSClientConfig, error) {
	tlsConfig := ssl.TLSClientConfig{}

	skipVerifyKey := fmt.Sprintf("%s.insecureSkipVerify", prefix)
	if val, err := String(skipVerifyKey); err == nil {
		skipVerifyVal := val
		if skipVerifyVal == "true" {
			tlsConfig.InsecureSkipVerify = true
		}
	}

	certFileKey := fmt.Sprintf("%s.certFile", prefix)
	if val, err := String(certFileKey); err == nil {
		tlsConfig.CertFile = val
	}

	keyFileKey := fmt.Sprintf("%s.keyFile", prefix)
	if val, err := String(keyFileKey); err == nil {
		tlsConfig.KeyFile = val
	}

	caFileKey := fmt.Sprintf("%s.caFile", prefix)
	if val, err := String(caFileKey); err == nil {
		tlsConfig.CAFile = val
	}

	passwordKey := fmt.Sprintf("%s.password", prefix)
	if val, err := String(passwordKey); err == nil {
		tlsConfig.Password = val
	}

	return tlsConfig, nil
}

// GetClientTLSConfig get client tls config
func GetClientTLSConfig(prefix string) (*tls.Config, error) {
	config, err := NewTLSClientConfigFromConfig(prefix)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}

	if len(config.CAFile) != 0 && len(config.CertFile) != 0 && len(config.KeyFile) != 0 {
		tlsConf, err = ssl.ClientTLSConfVerity(config.CAFile, config.CertFile, config.KeyFile, config.Password)
		if err != nil {
			return nil, err
		}
	}

	return tlsConf, nil
}
