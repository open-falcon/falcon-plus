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

package funcs

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

func UrlMetrics() (L []*model.MetricValue) {
	reportUrls := g.ReportUrls()
	sz := len(reportUrls)
	if sz == 0 {
		return
	}
	hostname, err := g.Hostname()
	if err != nil {
		hostname = "None"
	}
	for furl, timeout := range reportUrls {
		tags := fmt.Sprintf("url=%v,timeout=%v,src=%v", furl, timeout, hostname)
		if ok, _ := probeUrl(furl, timeout); !ok {
			L = append(L, GaugeValue(g.URL_CHECK_HEALTH, 0, tags))
			continue
		}
		L = append(L, GaugeValue(g.URL_CHECK_HEALTH, 1, tags))
	}
	return
}

func probeUrl(furl string, timeout string) (bool, error) {
	bs, err := sys.CmdOutBytes("curl", "--max-filesize", "102400", "-I", "-m", timeout, "-o", "/dev/null", "-s", "-w", "%{http_code}", furl)
	if err != nil {
		log.Printf("probe url [%v] failed.the err is: [%v]\n", furl, err)
		return false, err
	}
	reader := bufio.NewReader(bytes.NewBuffer(bs))
	retcode, err := file.ReadLine(reader)
	if err != nil {
		log.Println("read retcode failed.err is:", err)
		return false, err
	}
	if strings.TrimSpace(string(retcode)) != "200" {
		log.Printf("return code [%v] is not 200.query url is [%v]", string(retcode), furl)
		return false, err
	}
	return true, err
}
