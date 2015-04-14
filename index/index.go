package index

// 初始化索引功能模块
func StartIndex() {
	go StartIndexDeleteTask()
}
