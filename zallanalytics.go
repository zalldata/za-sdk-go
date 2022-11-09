/*
 * Created by sdk on 2020/01/06.
 * Copyright 2015－2020 Zall Data Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ZallAnalytics

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/zalldata/za-sdk-go/consumers"
	"github.com/zalldata/za-sdk-go/structs"
	"github.com/zalldata/za-sdk-go/utils"
)

const (
	TRACK             = "track"
	TRACK_SIGNUP      = "track_signup"
	PROFILE_SET       = "profile_set"
	PROFILE_SET_ONCE  = "profile_set_once"
	PROFILE_INCREMENT = "profile_increment"
	PROFILE_APPEND    = "profile_append"
	PROFILE_UNSET     = "profile_unset"
	PROFILE_DELETE    = "profile_delete"
	ITEM_SET          = "item_set"
	ITEM_DELETE       = "item_delete"

	BIND_DEVICE_ID   = 0
	BIND_LOGIN_ID    = 1
	BIND_MOBILE      = 2
	BIND_UNION_ID    = 3
	BIND_OPEN_ID     = 4
	BIND_EXTERNAL_ID = 5

	SDK_VERSION = "1.0.3"
	LIB_NAME    = "Golang"

	MAX_ID_LEN = 255
)

// 静态公共属性
var superProperties map[string]interface{}

type ZallAnalytics struct {
	C           consumers.Consumer
	ProjectName string
	TimeFree    bool
}

func InitZallAnalytics(c consumers.Consumer, projectName string, timeFree bool) ZallAnalytics {
	return ZallAnalytics{C: c, ProjectName: projectName, TimeFree: timeFree}
}

func (za *ZallAnalytics) track(etype, event, distinctId, originId string, properties map[string]interface{}, isLoginId bool) error {
	eventTime := utils.NowMs()
	if et := extractUserTime(properties); et > 0 {
		eventTime = et
	}

	data := structs.EventData{
		Type:          etype,
		Time:          eventTime,
		DistinctId:    distinctId,
		Properties:    properties,
		LibProperties: getLibProperties(),
	}

	if za.ProjectName != "" {
		data.Project = za.ProjectName
	}

	if etype == TRACK || etype == TRACK_SIGNUP {
		data.Event = event
	}

	if etype == TRACK_SIGNUP {
		data.OriginId = originId
	}

	if za.TimeFree {
		data.TimeFree = true
	}

	if isLoginId {
		properties["$is_login_id"] = true
	}

	err := data.NormalizeData()
	if err != nil {
		return err
	}

	return za.C.Send(data)
}

func (za *ZallAnalytics) Flush() {
	za.C.Flush()
}

func (za *ZallAnalytics) Close() {
	za.C.Close()
}

func (za *ZallAnalytics) Track(distinctId, event string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	// merge properties
	if properties == nil {
		nproperties = make(map[string]interface{})
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	// merge super properties
	if superProperties != nil {
		utils.MergeSuperProperty(superProperties, nproperties)
	}
	nproperties["$lib"] = LIB_NAME
	nproperties["$lib_version"] = SDK_VERSION

	return za.track(TRACK, event, distinctId, "", nproperties, isLoginId)
}

func (za *ZallAnalytics) TrackSignup(distinctId, originId string, distinctIdType, originIdType int) error {
	// check originId and merge properties
	if originId == "" {
		return errors.New("property [original_id] must not be empty")
	}
	if len(originId) > MAX_ID_LEN {
		return errors.New("the max length of property [original_id] is 255")
	}

	properties := make(map[string]interface{})
	// merge super properties
	if superProperties != nil {
		utils.MergeSuperProperty(superProperties, properties)
	}
	properties["$lib"] = LIB_NAME
	properties["$lib_version"] = SDK_VERSION

	properties["$distinctIdType"] = distinctIdType
	properties["$originalIdType"] = originIdType

	return za.track(TRACK_SIGNUP, "$SignUp", distinctId, originId, properties, false)
}

func (za *ZallAnalytics) ProfileSet(distinctId string, distinctIdType int, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}
	nproperties["$distinctIdType"] = distinctIdType
	return za.track(PROFILE_SET, "", distinctId, "", nproperties, isLoginId)
}

func (za *ZallAnalytics) ProfileSetOnce(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return za.track(PROFILE_SET_ONCE, "", distinctId, "", nproperties, isLoginId)
}

func (za *ZallAnalytics) ProfileIncrement(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return za.track(PROFILE_INCREMENT, "", distinctId, "", nproperties, isLoginId)
}

func (za *ZallAnalytics) ProfileAppend(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return za.track(PROFILE_APPEND, "", distinctId, "", nproperties, isLoginId)
}

func (za *ZallAnalytics) ProfileUnset(distinctId string, properties map[string]interface{}, isLoginId bool) error {
	var nproperties map[string]interface{}

	if properties == nil {
		return errors.New("property should not be nil")
	} else {
		nproperties = utils.DeepCopy(properties)
	}

	return za.track(PROFILE_UNSET, "", distinctId, "", nproperties, isLoginId)
}

func (za *ZallAnalytics) ProfileDelete(distinctId string, isLoginId bool) error {
	nproperties := make(map[string]interface{})

	return za.track(PROFILE_DELETE, "", distinctId, "", nproperties, isLoginId)
}

func (za *ZallAnalytics) ItemSet(itemType string, itemId string, properties map[string]interface{}) error {
	libProperties := getLibProperties()
	time := utils.NowMs()
	if properties == nil {
		properties = map[string]interface{}{}
	}

	itemData := structs.Item{
		Type:          ITEM_SET,
		ItemId:        itemId,
		Time:          time,
		ItemType:      itemType,
		Properties:    properties,
		LibProperties: libProperties,
	}

	err := itemData.NormalizeItem()
	if err != nil {
		return err
	}

	return za.C.ItemSend(itemData)
}

func (za *ZallAnalytics) ItemDelete(itemType string, itemId string) error {
	libProperties := getLibProperties()
	time := utils.NowMs()

	itemData := structs.Item{
		Type:          ITEM_DELETE,
		ItemId:        itemId,
		Time:          time,
		ItemType:      itemType,
		Properties:    map[string]interface{}{},
		LibProperties: libProperties,
	}

	err := itemData.NormalizeItem()
	if err != nil {
		return err
	}

	return za.C.ItemSend(itemData)
}

// 注册公共属性
func (za *ZallAnalytics) RegisterSuperProperties(superProperty map[string]interface{}) {
	if superProperties == nil {
		superProperties = make(map[string]interface{})
	}
	utils.MergeSuperProperty(superProperty, superProperties)
}

// 清除公共属性
func (za *ZallAnalytics) ClearSuperProperties() {
	superProperties = make(map[string]interface{})
}

// 清除指定 key 的公共属性
func (za *ZallAnalytics) UnregisterSuperProperty(key string) {
	delete(superProperties, key)
}

func getLibProperties() structs.LibProperties {
	lp := structs.LibProperties{}
	lp.Lib = LIB_NAME
	lp.LibVersion = SDK_VERSION
	lp.LibMethod = "code"
	if pc, file, line, ok := runtime.Caller(3); ok { //3 means sdk's caller
		f := runtime.FuncForPC(pc)
		lp.LibDetail = fmt.Sprintf("##%s##%s##%d", f.Name(), file, line)
	}

	return lp
}

func extractUserTime(p map[string]interface{}) int64 {
	if t, ok := p["$time"]; ok {
		v, ok := t.(int64)
		if !ok {
			fmt.Fprintln(os.Stderr, "It's not ok for type string")
			return 0
		}
		delete(p, "$time")

		return v
	}

	return 0
}

func InitDefaultConsumer(url string, timeout int) (*consumers.DefaultConsumer, error) {
	return consumers.InitDefaultConsumer(url, timeout)
}

func InitBatchConsumer(url string, max, timeout int) (*consumers.BatchConsumer, error) {
	return consumers.InitBatchConsumer(url, max, timeout)
}

func InitLoggingConsumer(filename string, hourRotate bool) (*consumers.LoggingConsumer, error) {
	return consumers.InitLoggingConsumer(filename, hourRotate)
}

func InitConcurrentLoggingConsumer(filename string, hourRotate bool) (*consumers.ConcurrentLoggingConsumer, error) {
	return consumers.InitConcurrentLoggingConsumer(filename, hourRotate)
}

func InitDebugConsumer(url string, writeData bool, timeout int) (*consumers.DebugConsumer, error) {
	return consumers.InitDebugConsumer(url, writeData, timeout)
}
