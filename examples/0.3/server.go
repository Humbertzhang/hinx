package main

import (
	"fmt"
	"hinx/hiface"
	"hinx/hnet"
)


// ping test 自定义路由
// 这个Router的功能：不论发送什么数据，我就给你写before ping   ping!ping!ping!   after ping ... 这几个字符串
type PingRouter struct {
	hnet.BaseRouter
}

// 重写几个继承过来的Handle，以实现实际功能
// PreHandle
func(pr *PingRouter) PreHandle (request hiface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n")); err != nil {
		fmt.Println("Call Router PreHandle Error:", err)
	}
}

// Handle
func(pr *PingRouter) Handle (request hiface.IRequest) {
	fmt.Println("Call Router Handle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("ping!ping!ping\n")); err != nil {
		fmt.Println("Call Router Handle Error:", err)
	}
}

// PostHandle
func(pr *PingRouter) PostHandle (request hiface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n")); err != nil {
		fmt.Println("Call Router PostHandle Error:", err)
	}
}

func main() {
	// 创建一个Server句柄
	s := hnet.NewServer("[hinx v0.3 ]")
	// 给Server注册Router
	s.AddRouter(&PingRouter{})
	// 启动
	s.Serve()
}
