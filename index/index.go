package index

import (
	"github.com/open-falcon/task/g"
	"log"
)

// 初始化索引功能模块
func Start() {
	if g.Config().Index.Enabled {
		StartDB()
		//StartIndexDeleteTask()
		StartIndexUpdateAllTask()
		log.Println("index:Start, ok")
	} else {
		log.Println("index:Start, index not enabled")
	}
}
