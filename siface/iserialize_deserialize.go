package siface

type ISerializeAndDeserialize interface {
	Serialize(data interface{}) ([]byte, error)

	Deserialize(bs []byte, data interface{}) error
}
