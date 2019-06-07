package hnet

import (
	"fmt"
	"hinx/hiface"
	"hinx/utils"
	"strconv"
)


// msgHandler 是处理客户端传过来的消息的handler
// 其需要做的是获得了request之后，将其分发给相应的处理函数
// 此外又负责了统筹任务队列与任务池的职责.
type MsgHandler struct {
	// 存放每一个MSGID 对应的处理方法
	Apis map[uint32] hiface.IRouter

	// 负责worker读取任务的消息队列
	// 每一个Chan即为一个消息队列，这里为一个消息队列的集合
	TaskQueue []chan hiface.IRequest
	// 业务工作Worker池worker数量
	WorkerPoolSize 	uint32
}


// 初始化
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]hiface.IRouter),
		// 每一个worker对应一个TaskQueue，因此数量都为workerPoolSize
		TaskQueue:      make([]chan hiface.IRequest, utils.GlobalObject.WorkerPoolSize),
		// 从全局配置参数中获取
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
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

// 启动一个工作池
// 开启工作池动作只发生一次，即一个zinx框架只能有一个工作池
func (mh *MsgHandler) StartWorkerPool() {
	fmt.Println("Worker Pool is starting...")
	// 开启PoolSize大小的worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 1. 给当前的worker对应的channel消息队列开辟空间

		mh.TaskQueue[i] = make(chan hiface.IRequest,  utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan hiface.IRequest) {
	fmt.Println("WorkerID = ", workerID, " is started...")
	for {
		// 用Select阻塞
		select {
		// 如果有消息传过来,处理这个Request
		case request := <- taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue，由Worker进行处理
func (mh *MsgHandler) SendMsgToQueue(request hiface.IRequest) {
	// 1.将消息平均分配给不同的worker，即轮训分配
	// 根据客户端建立的ConnID 进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
					"request MsgID =", request.GetMsgID(),
					" to WorkerID =", workerID)

	// 2.将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}