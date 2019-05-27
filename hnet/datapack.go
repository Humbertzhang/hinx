package hnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hinx/hiface"
	"hinx/utils"
)

type DataPack struct {

}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取头长
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen 4个字节
	// ID 4个字节
	// 共8字节
	return 8;
}

// 封包
// 将Message封装成一个消息
func (dp *DataPack) Pack(msg hiface.IMessage) ([]byte, error) {
	// 创建存放bytes字节流的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将Datalen写入buff
	// 采用小端法
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	// 将MessageID 写入buff
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID())
	if err != nil {
		return nil, err
	}

	// 将数据写入buff中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// 拆包
// 将bytes封装成IMessage
func (dp *DataPack) Unpack(binaryData []byte) (hiface.IMessage, error) {
	// 创建一个从二进制数据中读数据的IOReader
	reader := bytes.NewReader(binaryData)

	// 只解压head信息，得到len 和 ID
	msg := &Message{}

	// 读len
	err := binary.Read(reader, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 读ID
	err = binary.Read(reader, binary.LittleEndian, &msg.ID)
	if err != nil {
		return nil, err
	}

	// 限制包的大小
	if(utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("msg package too large")
	}
	// 此时msg中没有DATA
	return msg, nil
}