package client

import (
	"fmt"
	"os"
)

/**
 * 错误检查
 * @param err   错误类型数据
 * @param info  错误信息提示内容
 **/
func checkError(err error, info string) {

	if err != nil {

		fmt.Println(info + "  " + err.Error())

		os.Exit(1)

	}

}

/**
 * 聊天内容显示
 **/
func showMassage() {

	for {

		message := <-recvData
		fmt.Println(message)
	}

}

/**
 * 聊天数据输入
 **/
func inputMessage() {
	var input, objUser string
	for {

		fmt.Scanln(&input)

		switch input {

		// 用户退出
		case "/quit":
			fmt.Println("退出聊天室，欢迎下次使用！")
			os.Exit(0)

		// 向单个用户发送消息
		case "/to":
			fmt.Println("请输入聊天对象：")
			fmt.Scanln(&objUser)
			fmt.Println("请输入消息：")
			fmt.Scanln(&input)
			if len(input) != 0 {
				input = "To-" + uname + "-" + objUser + "-" + input
			}

			// 默认群发
		default:
			if len(input) != 0 {
				input = uname + " : " + input
			}
		}

		// 发送数据
		if len(input) != 0 {
			sentData <- input
			input = ""
		}
	}
}
