package snet

import (
	"fmt"
	"little-fast-mq/sd"
	"little-fast-mq/siface"
	"testing"
)

//ping test 自定义路由
type PingRouter struct {
	BaseRouter
}

//Test Handle
func (route *PingRouter) Handle(req siface.IRequest, res siface.IResponse) {

	var data interface{}
	err := sd.Jsd.Deserialize(req.GetData(), &data)
	if err != nil {
		fmt.Println("deserialize data error:", err)
		return
	}
	fmt.Println("request data :", data)

	var topic string
	err = sd.Jsd.Deserialize(req.GetReqTopic(), &topic)

	if err != nil {
		fmt.Println("deserialize topic error:", err)
		return
	}
	stu := Student{
		Name:   "sunpengwei",
		Age:    26,
		CardNo: "610525199202171334",
	}
	err = res.GetConnection().SendBuffMsg(topic, stu)

	if err != nil {
		fmt.Println("server writer error:", err)
	}
}

func TestServer_Serve(t *testing.T) {

	server := NewServer()

	server.AddRoute("test1", &PingRouter{})
	server.AddRoute("test2", &PingRouter{})
	server.AddRoute("test3", &PingRouter{})

	server.Serve()
}
