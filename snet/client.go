package snet

import (
	"errors"
	"fmt"
	"little-fast-mq/sd"
	"little-fast-mq/siface"
	"little-fast-mq/utils"
	"net"
)

type Client struct {
	Name     string
	Ip       string
	conn     *net.TCPConn
	isClosed bool
}

func (c *Client) Read() (interface{}, error) {

	dp := NewDataPack()

	headData := make([]byte, dp.GetHeadLen())
	_, err := c.conn.Read(headData)
	if err != nil {
		fmt.Println("client read head error:", err)
		return "", err
	}

	headMsg, err := dp.UnPack(headData)
	if err != nil {
		fmt.Println("client unpack error:", err)
		return "", err
	}

	if headMsg.GetTopicLen() <= 0 {
		return "", nil
	}
	topicbs := make([]byte, headMsg.GetTopicLen())

	_, err = c.conn.Read(topicbs)
	if err != nil {
		fmt.Println("client read body error:", err)
		return "", err
	}

	if headMsg.GetDataLen() <= 0 {
		return "", nil
	}
	databs := make([]byte, headMsg.GetDataLen())

	_, err = c.conn.Read(databs)
	if err != nil {
		fmt.Println("client read body error:", err)
		return "", err
	}
	//msg 是有data数据的，需要再次读取data数据
	msg := headMsg.(*Message)

	msg.SetData(databs)
	msg.SetTopic(topicbs)
	var t interface{}
	err = sd.Jsd.Deserialize(databs, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Client) Send(topic string, data interface{}) (int, error) {

	dp := NewDataPack()

	topicbs, err := sd.Jsd.Serialize(topic)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	databs, err := sd.Jsd.Serialize(data)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	msg := NewMessagePackage(utils.GetRandUint32(), topicbs, databs)

	w, err := dp.Pack(msg)

	if err != nil {
		fmt.Println("client write data error:", err)
		return 0, err
	}
	count, err := c.conn.Write(w)

	if err != nil {
		fmt.Println("client write data error:", err)
		return 0, err
	}
	return count, nil

}

func (c *Client) Stop() {

	if c.isClosed {
		return
	}
	c.isClosed = true
	err := c.conn.Close()
	if err != nil {

		fmt.Println("client closed fail", err)
		c.isClosed = false
		return
	}
	fmt.Println("client closed successful")

}
func NewClient() (siface.IClient, error) {

	locaIp := utils.GetLocalIp()
	if locaIp == "" {
		return nil, errors.New("get local ip is null")
	}

	raddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", locaIp, 7777))

	conn, err := net.DialTCP("tcp4", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Name:     "client",
		Ip:       locaIp,
		isClosed: false,
		conn:     conn,
	}, nil

}

/**
是否关闭
*/
func (c *Client) IsClosed() bool {

	return c.isClosed
}
