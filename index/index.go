package index

import (
	"log"
)

// 初始化索引功能模块
func Start() {
	StartIndexDeleteTask()
	log.Println("index:Start, ok")
}
