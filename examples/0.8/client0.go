package main

import (
	"fmt"
	"hinx/hnet"
	"io"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client0 started...Need Ping Router.")

	// 1. 连接服务器，得到conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	// 2. 尝试调用Write方法，写数据

	for {
		// 发送封包的Msg消息
		dp := hnet.NewDataPack()
		binaryMsg, err := dp.Pack(hnet.NewMsgPackage(0, []byte("Hello from client. Need Ping Router.")))
		if err != nil {
			fmt.Println("Pack msg error:", err)
			return
		}

		// 得到了封包后的数据后，直接写
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("Send msg error:", err)
			return
		}

		// 拆开服务器回复的数据

		// 先读取head部分
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("Get msg head error:", err)
			return
		}

		// 将二进制的head拆包到结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("unpack msg head error:", err)
			return
		}

		// msg中有数据
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*hnet.Message)
			// 再读取data
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("unpack msg data error:", err)
				return
			}
			fmt.Println("---> Receive Server Msg. ID = ", msg.GetMsgID(), " len = ", msg.GetMsgLen(),
							" data = ", string(msg.GetData()))
		}


		// cpu 阻塞，避免占用大量资源
		time.Sleep(1*time.Second)
	}

}
