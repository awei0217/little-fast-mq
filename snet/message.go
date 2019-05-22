package snet

type Message struct {
	Id       uint32 //消息的ID
	TopicLen uint32
	DataLen  uint32 //报文长度
	topic    []byte //主题
	Data     []byte //报文

}

func NewMessagePackage(id uint32, topic, data []byte) *Message {

	return &Message{
		Id:       id,
		TopicLen: uint32(len(topic)),
		DataLen:  uint32(len(data)),
		topic:    topic,
		Data:     data,
	}
}

func (m *Message) GetDataLen() uint32 {

	return m.DataLen
}

func (m *Message) GetData() []byte {

	return m.Data
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

//设计消息ID
func (msg *Message) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

//设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
	msg.SetDataLen(uint32(len(data)))
}

func (msg *Message) SetTopic(topic []byte) {
	msg.topic = topic
	msg.SetTopicLen(uint32(len(topic)))
}

func (msg *Message) GetTopic() []byte {
	return msg.topic
}

func (msg *Message) SetTopicLen(len uint32) {
	msg.TopicLen = len
}

func (msg *Message) GetTopicLen() uint32 {

	return msg.TopicLen
}
