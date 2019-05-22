package siface

type IServer interface {

	//启动服务
	Start()

	//停止服务
	Stop()

	//开启业务处理
	Serve()

	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRoute(topic string, router IRouter)

	GetConnManager() IConnManager

	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))

	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))

	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)

	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
