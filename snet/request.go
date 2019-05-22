package snet

import (
	"little-fast-mq/siface"
)

type Request struct {
	conn siface.IConnection //已经和客户端建立好的连接
	msg  siface.IMessage    //客户端请求的数据
}

func (req *Request) GetConnection() siface.IConnection {

	return req.conn
}

func (req *Request) GetData() []byte {

	return req.msg.GetData()
}

//获取请求的消息的ID
func (req *Request) GetReqMsgID() uint32 {

	return req.msg.GetMsgId()
}

//获取请求的消息的topic
func (req *Request) GetReqTopic() []byte {

	return req.msg.GetTopic()
}
