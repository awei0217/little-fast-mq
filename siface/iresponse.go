package siface

type IResponse interface {
	GetConnection() IConnection //获取请求连接信息

	GetData() []byte

	GetResMsgID() uint32 //获取响应消息的ID

}
