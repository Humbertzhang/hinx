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
	// 主要用于Reader告知Writer
	ExitChan chan bool

	// 用于Reader Writer之间通信的管道
	msgChan chan []byte

	// 消息管理msgID
	msgHandler hiface.IMsgHandler
}

// 初始化当前连接
func NewConnection(conn *net.TCPConn, connID uint32, handler hiface.IMsgHandler) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		isClosed: false,
		msgChan: make(chan []byte),
		ExitChan: make(chan bool, 1),
		msgHandler: handler,
	}
	return c
}


// 实现IConnection接口的方法
// 开启Reader Goroutine
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	// 结束之后调用Stop()
	defer fmt.Println(c.RemoteAddr().String(), "[Reader is exited]")
	defer c.Stop()

	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		// 拆包，得到msgID 和 msgLen，放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}

		// 根据dataLen 读取 data,存在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		// 将Conn封装为一个Request
		request := Request{
			conn: c,
			msg: msg,
		}

		// 从路由中找到对应的router
		// 根据绑定好的MsgID找到对应的api业务，执行
		go c.msgHandler.DoMsgHandler(&request)
	}

}

// 开启Writer Goroutine
// 用于写消息给客户端
// 将Reader和Writer职责分离
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exited.]")
	// 不断阻塞等待msgChan的消息
	for {
		select {
		// 有数据给客户端
		case data:= <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出
			return

		}
	}
}

// 启动连接:让当前的连接开始工作
func (c *Connection) Start() {
	fmt.Println("connection start...ConnID = ", c.ConnID)
	// 启动从当前连接 读数据 的Goroutine
	go c.StartReader()
	go c.StartWriter()
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

	// 将数据传到Writer
	c.msgChan <- msg
	return nil
}