package siface

type IMessage interface {
	GetDataLen() uint32

	GetMsgId() uint32

	GetData() []byte

	GetTopic() []byte

	GetTopicLen() uint32

	SetMsgId(id uint32)

	SetData(data []byte)

	SetDataLen(len uint32)

	SetTopic(topic []byte)

	SetTopicLen(len uint32)
}
