package hnet

import (
	"errors"
	"fmt"
	"hinx/hiface"
	"io"
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
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			c.ExitChan <- true
			continue
		}

		// 拆包，得到msgid 和 msglen，放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			c.ExitChan <- true
			continue
		}

		// 根据dataLen 读取 data,存在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				c.ExitChan <- true
				continue
			}
		}
		msg.SetData(data)

		// 将Conn封装为一个Request
		request := Request{
			conn: c,
			msg: msg,
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

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed before send msg")
	}

	// 将data封包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack error msgID = ", msgID)
		return errors.New("pack error msg")
	}

	// 写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msgID ", msgID, " error")
		c.ExitChan <- true
		return errors.New("conn Write error")
	}

	return nil
}