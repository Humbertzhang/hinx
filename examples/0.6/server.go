package main

import (
	"fmt"
	"hinx/hiface"
	"hinx/hnet"
)


//Hello router
// 用于演示多路由
type HelloRouter struct {
	hnet.BaseRouter
}

func (hr *HelloRouter) Handle(request hiface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	// 回写Hello
	err := request.GetConnection().SendMsg(201, []byte("Hello!"))
	if err != nil {
		fmt.Println("server error:", err)
	}
}

// ping 自定义路由
type PingRouter struct {
	hnet.BaseRouter
}

// Handle
func(pr *PingRouter) Handle (request hiface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
					 ", data = ", string(request.GetData()))

	// 现在都是回写Msg格式的数据，不再直接回写普通数据
	// 相当于对TCP做了一个简单的封装
	err := request.GetConnection().SendMsg(200, []byte("Ping Back!"))
	if err != nil {
		fmt.Println("server error:", err)
	}
}


func main() {
	// 创建一个Server句柄
	s := hnet.NewServer("[hinx v0.6]")

	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// 启动
	s.Serve()
}
