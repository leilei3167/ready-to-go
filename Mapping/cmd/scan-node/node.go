package main

import "mapping/internal/node"

func main() {
	//后续将以下全设为命令行参数
	//使用ulimit -a查看持有的文件句柄数量,并修改
	//sudo prlimit --nofile=65536 --pid $$; ulimit -n 65536 修改最大文件句柄

	node.NewNode("mongodb://localhost:27017", 40000,
		[]string{"124.223.174.63:9092", "182.61.6.67:9092"},
		"testGroup").Run()
}
