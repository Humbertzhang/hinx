package hiface

// 将请求的消息封装到Message中
// 处理数据层面的原子结构
type IMessage interface {
	//获取消息ID
	GetMsgID() uint32
	// 获取消息长度
	GetMsgLen() uint32
	// 获取消息数据
	GetData() []byte

	// 设置消息ID
	SetMsgID(uint32)
	// 设置消息数据
	SetData([]byte)
	// 设置消息长度
	SetMsgLen(uint32)
}

