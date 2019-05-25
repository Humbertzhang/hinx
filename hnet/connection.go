package hnet

import (
	"fmt"
	"hinx/hiface"
	"net"
)

// 连接模块
type Connection struct {
	// Socket
	Conn *net.TCPConn

	// 连接ID
	ConnID 	uint32

	// 当前连接状态
	isClosed bool

	// 告知当前连接已经退出的Channel
	ExitChan chan bool

	// 该连接的Router
	Router 	hiface.IRouter
}

// 初始化当前连接
func NewConnection(conn *net.TCPConn, connID uint32, router hiface.IRouter) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router: router,
	}
	return c
}


// 实现IConnection接口的方法

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running.")
	// 结束之后调用Stop()
	defer fmt.Println("connID = ", c.ConnID, " Reader is exited, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取数据到buf中
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error:", err)
			continue
		}

		// 将Conn封装为一个Request
		request := Request{
			conn: c,
			data: buf,
		}

		// 依次调用Connection注册的Router的PreHandle handle 和 PostHandle
		// 模板设计模式：先定义好了模板，再由使用者依次调用
		go func(request hiface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&request)
	}

}

// 启动连接:让当前的连接开始工作
func (c *Connection) Start() {
	fmt.Println("connection start...ConnID = ", c.ConnID)
	// 启动从当前连接 读数据 的Goroutine
	go c.StartReader()
	// TODO:启动从当前连接 写数据 的Goroutine
}
// 停止连接:结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("connection stop...ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}
// 获取当前连接绑定的Connection
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
// 获取远程客户端的TCP状态（IP Port）
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
// 发送数据，将数据发送给远程客户端
func (c *Connection) Send(data []byte) error {
	return nil
}