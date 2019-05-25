package hiface

// 路由Router: 将一个指令与其处理方式进行绑定的结构体
// 在这里是要将不同的消息对应不同的处理方式

type IRouter interface {
	// 处理conn业务之前的钩子方法Hook
	PreHandle (request IRequest)
	// 处理conn业务的主方法
	Handle (request IRequest)
	// 处理conn业务之后的钩子方法Hook
	PostHandle (request IRequest)
}