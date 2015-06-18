package store

import (
	"github.com/open-falcon/common/model"
	tlist "github.com/toolkits/container/list"
	tmap "github.com/toolkits/container/nmap"
)

const (
	defaultHistorySize = 3
)

var (
	HistoryCache = tmap.NewSafeMap()
)

func GetLastItem(key string) *model.GraphItem {
	itemlist, found := HistoryCache.Get(key)
	if !found || itemlist == nil {
		return &model.GraphItem{}
	}

	first := itemlist.(*tlist.SafeListLimited).Front()
	if first == nil {
		return &model.GraphItem{}
	}

	return first.(*model.GraphItem)
}

func GetAllItems(key string) []*model.GraphItem {
	ret := make([]*model.GraphItem, 0)
	itemlist, found := HistoryCache.Get(key)
	if !found || itemlist == nil {
		return ret
	}

	all := itemlist.(*tlist.SafeListLimited).FrontAll()
	for _, item := range all {
		if item == nil {
			continue
		}
		ret = append(ret, item.(*model.GraphItem))
	}
	return ret
}

func AddItem(key string, val *model.GraphItem) {
	itemlist, found := HistoryCache.Get(key)
	var slist *tlist.SafeListLimited
	if !found {
		slist = tlist.NewSafeListLimited(defaultHistorySize)
		HistoryCache.Put(key, slist)
	} else {
		slist = itemlist.(*tlist.SafeListLimited)
	}
	slist.PushFrontViolently(val)
}
