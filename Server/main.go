package main

import (
	"os"
	"xieyuschen/server"
)

func main() {
	if len(os.Args) != 2 {
		panic("输入错误！")
	}
	server.StartServer(os.Args[1])
}
