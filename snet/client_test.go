package snet

import (
	"fmt"
	"testing"
)

type Student struct {
	Name   string
	Age    int
	CardNo string
}

func TestNewClient(t *testing.T) {

	client, err := NewClient()

	if err != nil {
		fmt.Println("client connection remote server fail:", err)
		return
	}
	stu := Student{
		Name:   "sunpengwei",
		Age:    26,
		CardNo: "610525199202171334",
	}
	for {
		_, err := client.Send("test1", stu)

		if err != nil {
			break
		}
		data, err := client.Read()

		if err != nil {
			break
		}
		fmt.Println("response data:", data)

		break
	}

}
