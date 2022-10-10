package utils

import (
	"bytes"
	"sort"
	"strings"
	"sync"
)

var bufferPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

func SortedTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}
	size := len(tags)
	if size == 0 {
		return ""
	}
	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if size == 1 {
		for k, v := range tags {
			ret.WriteString(k)
			ret.WriteString("=")
			ret.WriteString(v)
		}
		return ret.String()
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for j, key := range keys {
		ret.WriteString(key)
		ret.WriteString("=")
		ret.WriteString(tags[key])
		if j != size-1 {
			ret.WriteString(",")
		}
	}

	return ret.String()
}

func PK(metric string, tagkeys []string) string {
	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)
	if len(tagkeys) == 0 {
		ret.WriteString(metric)
		return ret.String()
	}
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(strings.Join(tagkeys, ","))
	return ret.String()
}

func PKWithTagValue(metric string, tags map[string]string) string {
	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if len(tags) == 0 {
		ret.WriteString(metric)
		return ret.String()
	}
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(SortedTags(tags))
	return ret.String()
}

func SortedTagKVs(tags map[string]string) (tagKeys, tagValues []string) {
	if tags == nil {
		return
	}
	size := len(tags)
	if size == 0 {
		return
	}
	if size == 1 {
		for k, v := range tags {
			tagKeys = append(tagKeys, k)
			tagValues = append(tagValues, v)
		}
		return
	}
	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		tagKeys = append(tagKeys, key)
		tagValues = append(tagValues, tags[key])
	}
	return
}
