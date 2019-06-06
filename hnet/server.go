package hnet

import (
	"fmt"
	"hinx/hiface"
	"hinx/utils"
	"net"
)

// 实现IServer
type Server struct {
	Name 		string
	IPVersion 	string
	IP 			string
	Port 		int
	// MsgHandler 处理多个业务.
	msgHandler  hiface.IMsgHandler
}

func (s *Server) Start() {
	fmt.Printf("[Hinx] Server Name: %s, Start Server Listener at IP: %s, Port: %d.\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Hinx] Version: %s, MaxConn: %d, MaxPackageSize: %d.\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

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

			// 将s.Router与conn绑定
			// router首先在Server AddRouter时被设置，在处理Connection的生成Connection阶段
			// 被传给Connection
			// Connection在处理Request时进行调用
			dealConn := NewConnection(conn, cid, s.msgHandler)
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

func (s *Server) AddRouter(msgID uint32, router hiface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
	fmt.Println("Server add router success.")
}

//初始化Server模块
func NewServer(name string) hiface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.Port,
		msgHandler:NewMsgHandler(),
	}
	return s
}
