package index

import (
	"log"
)

// 初始化索引功能模块
func StartIndex() {
	go StartIndexDeleteTask()
	log.Println("StartIndex, ok")
}
