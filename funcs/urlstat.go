package funcs

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
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
		tags := fmt.Sprintf("url=%v,timeout=%v,src_ip=%v", furl, timeout, hostname)
		if ok, _ := probeUrl(furl, timeout); !ok {
			L = append(L, GaugeValue("url.check.health", 0, tags))
			continue
		}
		L = append(L, GaugeValue("url.check.health", 1, tags))
	}
	return
}

func probeUrl(furl string, timeout string) (bool, error) {
	bs, err := sys.CmdOutBytes("curl", "max-filesize", "102400", "-I", "-m", timeout, "-o", "/dev/null", "-s", "-w", "%{http_code}", furl)
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
