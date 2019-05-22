package siface

type IMsgHandler interface {
	DoMsgHandler(request IRequest) //马上已非阻塞的方式处理消息

	AddRouter(topic string, router IRouter) //为消息添加具体的处理逻辑

	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest) //将消息叫个队列，由worker处理

}
