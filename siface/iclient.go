package siface

type IClient interface {
	Stop()

	Read() (interface{}, error)

	Send(topic string, data interface{}) (int, error)
}
