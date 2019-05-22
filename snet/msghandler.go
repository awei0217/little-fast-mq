package snet

import (
	"fmt"
	"little-fast-mq/config"
	"little-fast-mq/sd"
	"little-fast-mq/siface"
)

type MsgHandle struct {
	Apis           map[string]siface.IRouter
	WorkerPoolSize uint32                 //业务工作worker池的数量
	TaskQueue      []chan siface.IRequest //worker负责取任务的消息的队列
}

func (mh *MsgHandle) DoMsgHandler(request siface.IRequest) {

	var topic string
	err := sd.Jsd.Deserialize(request.GetReqTopic(), &topic)

	if err != nil {
		fmt.Println("deserialize topic error:", err)
		return
	}
	handler, ok := mh.Apis[topic]

	if !ok {
		fmt.Println("api topic ", request.GetReqTopic(), " is not FOUND!")
		return
	}

	topicbs := []byte(request.GetReqTopic())
	response := &Response{
		conn: request.GetConnection(),
		msg:  NewMessagePackage(request.GetReqMsgID(), topicbs, nil),
	}
	handler.PostHandle(request)
	handler.Handle(request, response)
	handler.PostHandle(request)

}

func (mh *MsgHandle) AddRouter(topic string, router siface.IRouter) {
	mh.Apis[topic] = router
}

func (mh *MsgHandle) StartWorkerPool() {

	for i := 0; i < int(mh.WorkerPoolSize); i++ {

		mh.TaskQueue[i] = make(chan siface.IRequest, config.GlobalObject.MaxWorkerTaskLen)

		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan siface.IRequest) {

	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:

			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request siface.IRequest) {

	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则
	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}

func NewMsgHandle() siface.IMsgHandler {

	return &MsgHandle{
		Apis:           make(map[string]siface.IRouter),
		WorkerPoolSize: config.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan siface.IRequest, config.GlobalObject.WorkerPoolSize),
	}
}
