package hnet

import (
	"errors"
	"fmt"
	"hinx/hiface"
	"net"
)

// 实现IServer
type Server struct {
	Name 		string
	IPVersion 	string
	IP 			string
	Port 		int
}

func (s *Server) Start() {
	fmt.Printf("Start Server Listener At IP: %s, Port: %d.\n", s.IP, s.Port)

	// 将Start操作放到Goroutine中，避免阻塞
	go func() {
		// 1，创建一个TCP套接字，即一个Addr对象
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve TCP addr error.")
			return
		}

		// 2,监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen:", s.IPVersion, " err:", err)
			return
		}
		fmt.Println("Start Hinx Server success,", s.Name, "Listening...")

		var cid uint32
		cid = 0

		// 3,阻塞等待客户端链接，处理客户端链接业务
		for {
			// 客户端链接过来，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept Err:", err)
				continue
			}

			// 将该处理新连接的方法与conn绑定
			dealConn := NewConnection(conn, cid, CallbackToClient)
			cid++

			// 启动连接业务处理
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// TODO:将服务器资源状态或者一些已经开辟的链接信息进行停止或回收
}

// 用Serve()封装Start和Stop，用户仅需调用Serve即可
func (s *Server) Serve() {
	// 启动Server服务功能
	s.Start()

	// TODO: 做一些启动服务器之后的额外业务

	// 阻塞
	select {

	}

}


//初始化Server模块
func NewServer(name string) hiface.IServer {
	s := &Server{
		Name: name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8999,
	}
	return s
}

// Connection 回调函数
// 定义当前客户端连接的所绑定的handle
func CallbackToClient (conn *net.TCPConn, data []byte, cnt  int) error {
	// 回显业务
	fmt.Println("[Conn Handler] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back error:", err)
		return errors.New("CallbackToClient error")
	}
	return nil
}