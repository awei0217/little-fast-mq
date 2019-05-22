package siface

// 消息拆包粘包的接口

type IDataPack interface {
	GetHeadLen() uint32

	Pack(m IMessage) ([]byte, error) //粘包

	UnPack([]byte) (IMessage, error) //拆包
}
