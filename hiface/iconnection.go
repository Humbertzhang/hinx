package hiface

import "net"

// 定义连接模块
type IConnection interface {
	// 启动连接:让当前的连接开始工作
	Start()
	// 停止连接:结束当前连接的工作
	Stop()
	// 获取当前连接绑定的Connection
	GetTCPConnection() *net.TCPConn
	// 获取当前连接的ID
	GetConnID() uint32
	// 获取远程客户端的TCP状态（IP Port）
	RemoteAddr() net.Addr
	// 发送数据，将数据发送给远程客户端
	SendMsg(msgID uint32, data []byte) error
}


// 定义一个处理连接的方法
// *net.TCPConn: 连接
// []byte : 内容
// int: 内容长度
type HandleFunc func(*net.TCPConn, []byte, int) error