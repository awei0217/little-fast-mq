package snet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"little-fast-mq/config"
	"little-fast-mq/siface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {

	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {

	//Id uint32(4字节) +  DataLen uint32(4字节) + TopicLen uint32(4字节)
	return 12
}

func (dp *DataPack) Pack(msg siface.IMessage) ([]byte, error) {

	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//写topicLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetTopicLen()); err != nil {
		return nil, err
	}

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写topic数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetTopic()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法(解压数据)
func (dp *DataPack) UnPack(binaryData []byte) (siface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID和topicLen
	msg := &Message{}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//读topiclen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.TopicLen); err != nil {
		return nil, err
	}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if config.GlobalObject.MaxPacketSize > 0 && msg.DataLen > config.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
