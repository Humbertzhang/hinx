package hiface

// 服务器接口
type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行
	Serve()
}

