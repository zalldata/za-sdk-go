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
	c, err := sdk.InitDefaultConsumer("http://172.16.90.61:58080/a?service=zall&project=dddssss", 1000)
	if err != nil {
		fmt.Println(err)
		return
	}

	za := sdk.InitZallAnalytics(c, "default", false)
	defer za.Close()

	distinctId := "12345"
	event := "ViewProduct"
	properties := map[string]interface{}{
		"price":    12,
		"name":     "apple",
		"somedata": []string{"a", "b"},
	}

	err = za.Track(distinctId, event, properties, true)
	if err != nil {
		fmt.Println("track failed", err)
		return
	}

	fmt.Println("track done")
}
