package hiface

// 拆包 封包的方法

type IDataPack interface {
	// 获取头长
	GetHeadLen() uint32

	// 封包
	Pack(msg IMessage) ([]byte, error)

	// 拆包
	Unpack(data []byte) (IMessage, error)
}