package snet

import "little-fast-mq/siface"

type Response struct {
	conn siface.IConnection
	msg  siface.IMessage
}

func (res *Response) GetConnection() siface.IConnection {

	return res.conn
}

func (res *Response) GetData() []byte {

	return res.msg.GetData()
}

func (res *Response) GetResMsgID() uint32 {

	return res.msg.GetMsgId()
}
