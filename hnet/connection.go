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

	// 当前连接所绑定的处理业务的方法API
	handleAPI hiface.HandleFunc

	// 告知当前连接已经退出的Channel
	ExitChan chan bool
}

// 初始化当前连接
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi hiface.HandleFunc) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		isClosed: false,
		handleAPI: callbackApi,
		ExitChan: make(chan bool, 1),
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
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error:", err)
			continue
		}
		// 调用当前连接所绑定的handle
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID:", c.ConnID," Reader handler API error:", err)
			break
		}
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