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

package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func PostJsonBody(url string, v interface{}) (response []byte, err error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return response, err
	}

	bf := bytes.NewBuffer(bs)

	resp, err := http.Post(url, "application/json", bf)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)

}
