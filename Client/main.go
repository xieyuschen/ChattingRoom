/**
 * 基于Go语言的聊天室-客户端程序
 **/
package main

import (
	"fmt"
	"os"
	"xieyuschen/client"
)

func main() {
	if len(os.Args) == 2 {

		client.StartClient(os.Args[1]) // 启动客户端

	} else {

		fmt.Println("输入错误！")

	}
}
