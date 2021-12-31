package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

func Receive(recvData chan string, conn net.Conn) {

	for {

		buf := make([]byte, 1024)

		n, err := conn.Read(buf)

		if err != nil {

			connection <- conn

			return
		}

		recvData <- string(buf[:n])

	}

}

func sent(sentData chan string, conn net.Conn) {

	for {
		data := <-sentData
		_, err := conn.Write([]byte(data))
		if err != nil {
			connection <- conn
			return
		}
	}
}

func connManager(connection chan net.Conn) {
	for {
		conn := <-connection
		username := ipToUname[conn.RemoteAddr().String()]
		notesInfo(messages, "用户:"+username+"已退出聊天室！")
		conn.Close()
		delete(connections, username)
		delete(ipToUname, conn.RemoteAddr().String())
	}
}

func dataSnt(conns *map[string]clientInfo, messages chan string) {

	for {
		msg := <-messages

		fmt.Println(msg)

		data := strings.Split(msg, "-") // 聊天数据分析:
		length := len(data)

		if length == 2 { // 管理员单个用户发送控制命令

			(*conns)[data[0]].sentData <- data[1]

		} else if length == 3 { // 用户列表

			(*conns)[data[1]].sentData <- data[2]

		} else if length == 4 { // 向单个用户发送数据

			msg = data[1] + " say to you : " + data[3]

			(*conns)[data[2]].sentData <- msg

		} else {
			// 群发
			for _, value := range *conns {
				value.sentData <- msg
			}
		}

	}
}

/**
 * 用户登录之后将用户接受的数据放入公共channel
 * @param messages 数据通道
 **/
func dataRcv(conns *map[string]clientInfo, username string, recvData, messages chan string) {
	for {
		data := <-recvData
		if _, ok := (*conns)[username]; !ok {
			return
		}
		if len(data) > 0 {
			messages <- data
		}

	}

}

/**
 * 设置TCP连接
 * @param  tcpAddr  TCP地址格式
 * return net.TCPListener
 **/
func createTcp(addr string) net.Listener {
	path := "D:\\xieyuschen\\ChattingRoom\\Cert\\"
	cert, err := tls.LoadX509KeyPair(path+"server.pem", path+"server.key")
	if err != nil {
		panic("Failed to create tls listening," + err.Error())
	}
	certBytes, err := ioutil.ReadFile(path + "client.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	l, err := tls.Listen("tcp", addr, config)
	errorExit(err, "DialTCP")
	return l
}
