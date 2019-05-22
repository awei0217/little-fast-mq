package snet

import (
	"errors"
	"fmt"
	"io"
	"little-fast-mq/config"
	"little-fast-mq/sd"
	"little-fast-mq/siface"
	"little-fast-mq/utils"
	"net"
	"sync"
)

type Connection struct {
	TcpServer    siface.IServer //当前conn属于那个server
	Conn         *net.TCPConn
	ConnID       uint32         //该连接的id
	isClosed     bool           //该连接是否关闭
	Router       siface.IRouter //该连接的处理方法router
	ExitBuffChan chan bool      //告知该连接已经推出的chan
	MsgHandler   siface.IMsgHandler
	msgChan      chan []byte //无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan  chan []byte //有关冲管道，用于读、写两个goroutine之间的消息通信

	// ================================
	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
	// ================================

}

//创建连接的方法
func NewConnection(server siface.IServer, conn *net.TCPConn, connID uint32, msgHandler siface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, config.GlobalObject.MaxMsgChanLen), //不要忘记初始化
		property:     make(map[string]interface{}),                         //对链接属性map初始化

	}

	//将创建的连接添加到连接管理中
	c.TcpServer.GetConnManager().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		// 创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		//拆包，得到msgid 和 datalen 和 topiclen放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		//根据topiclen 读取 topic放入msg中
		if msg.GetTopicLen() > 0 {
			topic := make([]byte, msg.GetTopicLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), topic); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
			msg.SetTopic(topic)
		}

		//根据 dataLen 读取 data，放在msg.Data中
		if msg.GetDataLen() > 0 {
			data := make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
			msg.SetData(data)
		}

		//得到当前客户端请求的Request数据
		req := &Request{
			conn: c,
			msg:  msg, //将之前的buf 改成 msg
		}
		if config.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {
	defer fmt.Println(c.RemoteAddr().String(), " conn Writer exit!")

	for {

		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:", err, " conn Writer exit")
				return
			}
			//针对有缓冲channel需要些的数据处理
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				break
				fmt.Println("msgBuffChan is Closed")
			}
		case <-c.ExitBuffChan:
			//通道关闭
			return

		}
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {

	//开启从客户端读取数据流程的Goroutine
	go c.StartReader()

	//开启从往客户端写数据流程的Goroutine
	go c.StartWriter()

	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	c.TcpServer.GetConnManager().Remove(c)
	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(topic string, data interface{}) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	topicbs, err := sd.Jsd.Serialize(topic)
	if err != nil {
		return errors.New(err.Error())
	}
	databs, err := sd.Jsd.Serialize(data)
	if err != nil {
		return errors.New(err.Error())
	}
	msg, err := dp.Pack(NewMessagePackage(utils.GetRandUint32(), topicbs, databs))
	if err != nil {
		fmt.Println("Pack error msg id = ", utils.GetRandUint32())
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBuffMsg(topic string, data interface{}) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()

	topicbs, err := sd.Jsd.Serialize(topic)
	if err != nil {
		return errors.New(err.Error())
	}
	databs, err := sd.Jsd.Serialize(data)
	if err != nil {
		return errors.New(err.Error())
	}

	msg, err := dp.Pack(NewMessagePackage(utils.GetRandUint32(), topicbs, databs))
	if err != nil {
		fmt.Println("Pack error msg id = ", utils.GetRandUint32())
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
