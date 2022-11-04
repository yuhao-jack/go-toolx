// Package netx
// @Description: 网络工具包
package netx

import (
	"io"
	"log"
	"net"
)

const (
	KEEPALIVE     = "KEEP_ALIVE"
	NEWCONNECTION = "NEW_CONNECTION"
)

// CreateTCPListener
//
//	@Description: 创建tcp监听对象
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-02 11:06:23
//	@param addr 要监听的地址
//	@return *net.TCPListener tcp监听对象
//	@return error
func CreateTCPListener(addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpListener, nil
}

// CreateTcpConn
//
//	@Description:  创建tcp链接
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-02 11:07:13
//	@param addr
//	@return *net.TCPConn
//	@return error
func CreateTcpConn(addr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}

// Join2Conn
//
//	@Description: 交换两个链接的IO
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-02 11:07:29
//	@param local
//	@param remote
func Join2Conn(local *net.TCPConn, remote *net.TCPConn) {
	go joinConn(local, remote)
	go joinConn(remote, local)
}

// joinConn
//
//	@Description: 交互io
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-02 11:08:03
//	@param local
//	@param remote
func joinConn(local *net.TCPConn, remote *net.TCPConn) {
	defer local.Close()
	defer remote.Close()
	_, err := io.Copy(local, remote)
	if err != nil {
		log.Default().Println("copy failed ", err.Error())
		return
	}
}
