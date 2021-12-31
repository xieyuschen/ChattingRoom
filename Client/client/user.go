package client

import "fmt"

var uname string

/**
 * 登录程序
 **/
func clientLogin() {

LOOP:
	for {
		var username string
		fmt.Println("请输入用户名：")
		fmt.Scanln(&username)
		if username == "" {
			goto LOOP
		}

		sentData <- "Log-" + username
		res := <-recvData
		if res == "success" {
			fmt.Println("登录成功，你已进入聊天室！")
			uname = username
			go showMassage()
			break
		} else {

			fmt.Println(res)
		}
	}

	inputMessage() // 启动聊天内容输入

}
