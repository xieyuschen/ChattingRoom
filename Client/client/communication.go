package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
)

var conn net.Conn

var recvData = make(chan string)
var sentData = make(chan string)

/**
 * 创建TCP连接
 * @param  tcpAddr  TCP地址格式
 * return net.Conn
 **/
func createTCP(tcpaddr string) net.Conn {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)

	checkError(err, "ResolveTCPAddr")

	path := "D:\\xieyuschen\\ChattingRoom\\Cert\\"

	cert, err := tls.LoadX509KeyPair(path+"client.pem", path+"client.key")
	if err != nil {
		panic(err)
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
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", tcpAddr.String(), config)

	checkError(err, "DialTCP")

	fmt.Print("连接服务器成功，")

	return conn
}

func StartClient(tcpAddr string) {

	conn = createTCP(tcpAddr) // 创建TCP连接

	go recv(recvData, conn) // 启动数据接收线程

	go sent(sentData, conn) // 启动数据发送线程
	clientLogin()           // 登录
}

/**
 * 数据接收
 * @param  recvData  接收数据Channel
 * @param  conn      TCP连接对象指针
 **/
func recv(recvData chan string, conn net.Conn) {

	for {

		buf := make([]byte, 1024)

		n, err := conn.Read(buf)

		checkError(err, "Connection")

		recvData <- string(buf[:n])

	}

}

/**
 * 数据发送
 * @param  sentData  接收数据Channel
 * @param  conn      TCP连接对象
 **/
func sent(sentData chan string, conn net.Conn) {

	for {

		data := <-sentData

		_, err := conn.Write([]byte(data))

		checkError(err, "Connection")

	}
}
