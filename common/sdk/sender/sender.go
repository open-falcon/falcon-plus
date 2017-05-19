// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sender

import (
	"log"
	"time"
)

const LIMIT = 200

var MetaDataQueue = NewSafeLinkedList()
var PostPushUrl string
var Debug bool

func StartSender() {
	go startSender()
}

func startSender() {
	for {
		L := MetaDataQueue.PopBack(LIMIT)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}

		err := PostPush(L)
		if err != nil {
			log.Println("[E] push to transfer fail", err)
		}
	}
}
