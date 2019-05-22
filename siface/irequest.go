package siface

type IRequest interface {
	GetConnection() IConnection //获取请求连接信息

	GetData() []byte //获取请求数据信息

	GetReqMsgID() uint32 //获取请求消息的ID

	GetReqTopic() []byte
}
