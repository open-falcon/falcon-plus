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

package store

import (
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	tlist "github.com/toolkits/container/list"
	tmap "github.com/toolkits/container/nmap"
)

const (
	defaultHistorySize = 3
)

var (
	// mem:  front = = = back
	// time: latest ...  old
	HistoryCache = tmap.NewSafeMap()
)

func GetLastItem(key string) *cmodel.GraphItem {
	itemlist, found := HistoryCache.Get(key)
	if !found || itemlist == nil {
		return &cmodel.GraphItem{}
	}

	first := itemlist.(*tlist.SafeListLimited).Front()
	if first == nil {
		return &cmodel.GraphItem{}
	}

	return first.(*cmodel.GraphItem)
}

func GetAllItems(key string) []*cmodel.GraphItem {
	ret := make([]*cmodel.GraphItem, 0)
	itemlist, found := HistoryCache.Get(key)
	if !found || itemlist == nil {
		return ret
	}

	all := itemlist.(*tlist.SafeListLimited).FrontAll()
	for _, item := range all {
		if item == nil {
			continue
		}
		ret = append(ret, item.(*cmodel.GraphItem))
	}
	return ret
}

func AddItem(key string, val *cmodel.GraphItem) {
	itemlist, found := HistoryCache.Get(key)
	var slist *tlist.SafeListLimited
	if !found {
		slist = tlist.NewSafeListLimited(defaultHistorySize)
		HistoryCache.Put(key, slist)
	} else {
		slist = itemlist.(*tlist.SafeListLimited)
	}

	// old item should be drop
	first := slist.Front()
	if first == nil || first.(*cmodel.GraphItem).Timestamp < val.Timestamp { // first item or latest one
		slist.PushFrontViolently(val)
	}
}
