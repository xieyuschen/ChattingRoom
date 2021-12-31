/**
 * 基于Go语言的聊天室-服务端程序
 **/
package server

import (
	"fmt"
	"net"
)

// 常量配置
const (
	dataFileName = "database.txt" // 用户注册数据文件保存名
	dataSize     = 5 * 1024       // 文件大小
	maxLog       = 30             // 最大同时登录人数
	maxReg       = 500            // 最大用户注册人数
	success      = "success"      // 登录注册返回给客户端成功标识
)

// 用户数据结构体
type userData struct {
	Password string
	Level    int
}

// 声明客户端的结构体
type clientInfo struct {
	conn     net.Conn    // 客户端的TCP连接对象
	sentData chan string // 服务器向客户端发送数据通道

}

var connections = make(map[string]clientInfo) // 声明成功登录之后的连接对象map

var uData = make(map[string]userData) // 声明登录用户数据map

var messages = make(chan string) // 声明消息channel

var adminList = make(map[string]string) // 声明管理员列表

var connection = make(chan net.Conn) // 声明连接管理Channel

var ipToUname = make(map[string]string) // ip地址与用户名队名

/**
 * 发送系统通知消息
 * @param messages  channel
 * @param info   channel
 **/
func notesInfo(messages chan string, info string) {
	messages <- info
}

func StartServer(port string) {
	l := createTcp(":" + port)
	fmt.Println("服务端启动成功，正在监听端口！")

	go dataSnt(&connections, messages)
	go connManager(connection)

	for {
		conn, err := l.Accept()
		if checkError(err, "Accept") == false {
			continue
		}
		fmt.Println("客户端:", conn.RemoteAddr().String(), "连接服务器成功！")
		var recvData = make(chan string)
		var sentData = make(chan string)
		go Receive(recvData, conn)             // 开启对客户端的接受数据线程
		go sent(sentData, conn)                // 开启对客户端的发送数据线程
		go userAuth(&conn, recvData, sentData) // 用户资格认证
	}
}
