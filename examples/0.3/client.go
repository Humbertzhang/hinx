package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟ECHO客户端
func main() {
	fmt.Println("client started...")

	// 1. 连接服务器，得到conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	// 2. 尝试调用Write方法，写数据

	for {
		_, err := conn.Write([]byte("Hello, Hinx V0.1"))
		if err != nil {
			fmt.Println("write conn error:", err)
			return
		}

		buf := make([]byte, 512)
		// 从服务端回写处获取输出
		count, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error:", err)
			return
		}
		fmt.Printf("server echo: %s, cnt = %d\n", buf, count)

		// cpu 阻塞，避免占用大量资源
		time.Sleep(1*time.Second)
	}

}
