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

package main

import (
	"fmt"
	sdk "github.com/zalldata/za-sdk-go"
)

func main() {
	//https://logcollect.zalldata.cn/a?project=iccvigrt&service=fl&token=21c6761a470d3a01a1e1c3ffcfcec751
	c, err := sdk.InitBatchConsumer("http://logcollect.zalldata.cn/a?service=zall&project=z7adds", 3, 1000)

	//c, err := sdk.InitBatchConsumer("http://172.16.90.61:58080/a?service=zall&project=dddssss", 3, 1000)
	if err != nil {
		fmt.Println(err)
		return
	}

	za := sdk.InitZallAnalytics(c, "z7adds", false)
	defer za.Close()

	distinctId := "ABCDEF123456"
	openId := "openId"
	unionId := "unionId"

	event := "ViewProduct"
	properties := map[string]interface{}{
		"$ip":            "2.2.2.2",
		"ProductId":      "1234562",
		"ProductCatalog": "Laptop Computer",
		"IzaddedToFav":   true,
	}

	err = za.Track(unionId, event, properties, false)
	if err != nil {
		fmt.Println("track failed", err)
		return
	}

	err = za.TrackSignup(unionId, openId, sdk.BIND_UNION_ID, sdk.BIND_OPEN_ID)
	if err != nil {
		fmt.Println("track failed", err)
		return
	}

	err = za.ProfileSet(distinctId, sdk.BIND_LOGIN_ID, properties, true)
	if err != nil {
		fmt.Println("track failed", err)
		return
	}

	fmt.Println("track done")
}
