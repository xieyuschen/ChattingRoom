package server

import (
	"fmt"
	"net"
	"strings"
)

/**
 * 管理员登录、注册
 **/
func intoManager() {
	fmt.Println("请输入将要进行操作：1、管理员注册 2、管理员登录")
	var input string
LOOP:
	{
		fmt.Scanln(&input)
		switch input {
		case "1":
			adminReg()
		case "2":
			adminLog()
		default:
			goto LOOP
		}
	}
	admimManager(messages)
}

/**
 * 管理员登录程序
 **/
func adminLog() {

	for {
		var adminname, password string
		fmt.Println("请输入管理员用户名：")
		fmt.Scanln(&adminname)
		fmt.Println("请输入管理员密码：")
		fmt.Scanln(&password)
		if pwd, ok := adminList[adminname]; !ok {
			fmt.Println("用户名或者密码错误")
		} else {
			if pwd != password {
				fmt.Println("用户名或者密码错误！")
			} else {
				fmt.Println("登录成功！")
				break
			}
		}
	}

}

/**
 * 管理员注册程序
 **/
func adminReg() {
	var adminname, password string
	fmt.Println("请输入管理员用户名：")
	fmt.Scanln(&adminname)
	fmt.Println("请输入管理员密码：")
	fmt.Scanln(&password)
	adminList[adminname] = password
	fmt.Println("注册成功！请登录")
	adminLog()
}

/**
 * 管理员管理模块
 **/
func admimManager(messages chan string) {

	for {
		var input, objUser string
		fmt.Scanln(&input)
		switch input {
		case "/to":
			fmt.Println(userList())
			fmt.Println("请输入聊天对象：")
			fmt.Scanln(&objUser)
			if _, ok := connections[objUser]; !ok {
				fmt.Println("不存在此用户!")
			} else {
				fmt.Println("请输入消息：")
				fmt.Scanln(&input)
				notesInfo(messages, "To-Manager-"+objUser+"-"+input)
			}

		case "/all":
			fmt.Println("请输入消息：")
			fmt.Scanln(&input)
			notesInfo(messages, "Manager say : "+input)

		case "/shield":
			fmt.Println(userList())
			fmt.Println("请输入屏蔽用户名：")
			fmt.Scanln(&objUser)
			notesInfo(messages, objUser+"-/shield")
			notesInfo(messages, "用户："+objUser+"已被管理员禁言！")

		case "/remove":
			fmt.Println(userList())
			fmt.Println("请输入踢出用户名：")
			fmt.Scanln(&objUser)
			notesInfo(messages, "用户："+objUser+"已被管理员踢出聊天室！")
			if _, ok := connections[objUser]; !ok {
				fmt.Println("不存在此用户!")
			} else {
				connections[objUser].conn.Close() // 删除该用户的连接
				delete(connections, objUser)      // 从已登录的列表中删除该用户
			}
		}

	}
}

/**
 * 获取用户列表
 * @param string   用户列表
 **/
func userList() string {
	var userList string = "当前在线用户列表："
	for user := range connections {

		userList += "\r\n" + user

	}
	return userList
}

func joinIn(username string) (status bool, info string) {

	if len(connections) == maxLog {
		return false, "当前登录人数已满,请稍后登录"
	}
	if _, ok := uData[username]; ok {
		return false, "用户已存在，请更换注册名!"
	} else {
		uData[username] = userData{}
		return true, success
	}
}

func userAuth(conn *net.Conn, recvData, sentData chan string) {

	for {
		data := strings.Split(<-recvData, "-") // 等待用户发送登录或注册数据
		username := data[1]
		status, info := joinIn(username)
		sentData <- info
		if status == true {
			messages <- ("用户:" + username + "进入聊天室")
			connections[username] = clientInfo{*conn, sentData}
			ipToUname[(*conn).RemoteAddr().String()] = username
			messages <- "List-" + username + "-" + userList()
			go dataRcv(&connections, username, recvData, messages)
			break
		}
	}
}
