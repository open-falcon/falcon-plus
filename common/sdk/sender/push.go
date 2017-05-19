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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/open-falcon/falcon-plus/common/model"
)

func PostPush(L []*model.JsonMetaData) error {
	bs, err := json.Marshal(L)
	if err != nil {
		return err
	}

	bf := bytes.NewBuffer(bs)

	resp, err := http.Post(PostPushUrl, "application/json", bf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	content := string(body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d != 200, response: %s", resp.StatusCode, content)
	}

	if Debug {
		log.Println("[D] response:", content)
	}

	return nil
}
