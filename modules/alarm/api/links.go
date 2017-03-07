package api

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/toolkits/net/httplib"
	"time"
)

func LinkToSMS(content string) (string, error) {
	links := g.Config().Api.Links
	uri := fmt.Sprintf("%s/store", links)
	req := httplib.Post(uri).SetTimeout(3*time.Second, 10*time.Second)
	req.Body([]byte(content))
	return req.String()
}
