package sd

import "encoding/json"

type JsonSerializerAndDeserialize struct {
}

var Jsd *JsonSerializerAndDeserialize

func (jsd *JsonSerializerAndDeserialize) Serialize(data interface{}) ([]byte, error) {

	return json.Marshal(data)
}

func (jsd *JsonSerializerAndDeserialize) Deserialize(bs []byte, data interface{}) error {

	return json.Unmarshal(bs, data)
}
func init() {

	Jsd = &JsonSerializerAndDeserialize{}
}
