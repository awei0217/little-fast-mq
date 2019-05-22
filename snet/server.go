//speed net
package snet

import (
	"fmt"
	"little-fast-mq/config"
	"little-fast-mq/siface"
	"net"
)

type Server struct {
	Name        string //服务器名称
	IPVersion   string //tcp4 or other
	IP          string //ip地址
	Port        int    //端口号
	msgHandle   siface.IMsgHandler
	ConnManager siface.IConnManager

	// =======================
	//新增两个hook函数原型

	//该Server的连接创建时Hook函数
	OnConnStart func(conn siface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn siface.IConnection)

	// =======================
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[little-fast-mq] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		config.GlobalObject.Version,
		config.GlobalObject.MaxConn,
		config.GlobalObject.MaxPacketSize)

	go func() {
		//开启业务处理工作池
		s.msgHandle.StartWorkerPool()
		// 启动服务
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			fmt.Println("resolve tcp addr err", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			fmt.Println("listener tcp err", err)
			return
		}
		fmt.Println("start speed-mq successful ", s.Name, " listener....")

		var connID uint32
		connID = 0

		for {
			conn, err := listener.AcceptTCP()

			if err != nil {
				fmt.Println("accept tcp err", err)
				continue
			}

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.GetConnManager().Len() >= config.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, connID, s.msgHandle)
			connID++
			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}

	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] speed-mq server , name ", s.Name)

	//清理所有连接
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	//阻塞,否则主Go退出， listener的go将会退出
	select {}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRoute(topic string, router siface.IRouter) {

	s.msgHandle.AddRouter(topic, router)

	fmt.Println("Add Router successful! ")
}

func (s *Server) GetConnManager() siface.IConnManager {

	return s.ConnManager
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(siface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(siface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn siface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn siface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

/*
  创建一个服务器句柄
*/
func NewServer() siface.IServer {
	//先初始化全局配置文件
	config.GlobalObject.Reload()

	s := &Server{
		Name:        config.GlobalObject.Name, //从全局参数获取
		IPVersion:   "tcp4",
		IP:          config.GlobalObject.Host,    //从全局参数获取
		Port:        config.GlobalObject.TcpPort, //从全局参数获取
		msgHandle:   NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}
