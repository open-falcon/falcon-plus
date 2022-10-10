package rpc

import (
	"regexp"
	"sync"
	"unicode"

	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/g"
	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/utils"
	log "github.com/sirupsen/logrus"
)

const MetricNameChars = "[^a-zA-Z0-9_:]"
const TagKeyNameChars = "[^a-zA-Z0-9_]"

var metricNameRegexp, tagKeyNameCharsRegexp *regexp.Regexp

func init() {
	metricNameRegexp, _ = regexp.Compile(MetricNameChars)
	tagKeyNameCharsRegexp, _ = regexp.Compile(TagKeyNameChars)
}

type P8sRelay int

func (this *P8sRelay) Ping(req NullRpcRequest, resp *SimpleRpcResponse) error {
	return nil
}

// Metric转换规范：[^a-zA-Z0-9_:] 会被替换为下划线，数字开头会添加 F_ 前缀
// Tag key转换规范：[^a-zA-Z0-9_] 会被替换为下划线，数字开头会添加 F_ 前缀
func (this *P8sRelay) Send(items []*P8sItem, resp *SimpleRpcResponse) error {
	wg := sync.WaitGroup{}
	for _, item := range items {
		wg.Add(1)
		go func(item *P8sItem) {
			defer wg.Done()
			if item.MetricType == g.Counter && item.Value < 0 {
				log.Printf("Error counter value is less than zero. %v\n", *item)
				return
			}
			item.Metric = metricNameRegexp.ReplaceAllString(item.Metric, "_")
			if unicode.IsDigit(rune(item.Metric[0])) {
				item.Metric = "F_" + item.Metric
			}
			item.Tags["hostname"] = item.Endpoint
			for key, value := range item.Tags {
				delete(item.Tags, key)
				key = tagKeyNameCharsRegexp.ReplaceAllString(key, "_")
				if unicode.IsDigit(rune(key[0])) {
					key = "F_" + key
				}
				item.Tags[key] = value
			}
			item.TagKeys, item.TagValues = utils.SortedTagKVs(item.Tags)
			item.PK = utils.PK(item.Metric, item.TagKeys)
			item.PKWithTagValue = utils.PKWithTagValue(item.Metric, item.Tags)
			g.P8sItemQueue.Append(item)
		}(item)
	}
	wg.Wait()
	return nil
}
