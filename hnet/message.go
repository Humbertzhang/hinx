package hnet


type Message struct {
	ID 				uint32 		// 消息ID
	DataLen 		uint32 		// 消息长度
	Data 			[]byte		// 消息
}

func NewMsgPackage(msgID uint32, data []byte) *Message {
	msg := &Message{
		ID:      msgID,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	return msg
}

//获取消息ID
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

// 获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息数据
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetMsgID(ID uint32) {
	m.ID = ID
}

// 设置消息数据
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// 设置消息长度
func (m *Message) SetMsgLen(length uint32) {
	m.DataLen = length
}