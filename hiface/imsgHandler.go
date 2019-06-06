package hiface

//消息管理抽象层
//负责将客户端消息和对应的Router对应起来

type IMsgHandler interface {
	// 调度对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// 为消息添加对应的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}