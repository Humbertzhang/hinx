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

// Handle
func(pr *PingRouter) Handle (request hiface.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
					 ", data = ", string(request.GetData()))

	// 现在都是回写Msg格式的数据，不再直接回写普通数据
	// 相当于对TCP做了一个简单的封装
	err := request.GetConnection().SendMsg(1, []byte("ping back wowowowowowowo."))
	if err != nil {
		fmt.Println("server error:", err)
	}
}


func main() {
	// 创建一个Server句柄
	s := hnet.NewServer("[hinx v0.5]")
	// 给Server注册Router
	s.AddRouter(&PingRouter{})
	// 启动
	s.Serve()
}
