package hiface

// irequest 接口
// 将客户端的 链接信息 和 请求的数据 包装到一个Request中
// 让Request作为客户端每次请求的原子数据结构


type IRequest interface {
	// 得到当前链接
	GetConnection() IConnection

	// 得到Request中的消息数据
	GetData() []byte

	// 获取请求的消息ID
	GetMsgID() uint32
}
