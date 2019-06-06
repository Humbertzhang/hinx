package hnet

import (
	"fmt"
	"hinx/hiface"
	"strconv"
)

type MsgHandler struct {
	// 存放每一个MSGID 对应的处理方法
	Apis map[uint32] hiface.IRouter
}


// 初始化
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]hiface.IRouter),
	}
}

// 调度
func (mh *MsgHandler) DoMsgHandler(request hiface.IRequest) {
	// 从Request中获取MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " not found. need register.")
	}
	// 根据MsgID 调用对应业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 添加
func (mh *MsgHandler) AddRouter(msgID uint32, router hiface.IRouter) {
	// 判断当前ID有无被注册
	// ok 代表已经注册过了
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 直接添加
	mh.Apis[msgID] = router
	fmt.Println("msgID = ", msgID, " success.")
}