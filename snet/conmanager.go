package snet

import (
	"errors"
	"fmt"
	"little-fast-mq/siface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]siface.IConnection //存储连接的集合

	connLock sync.RWMutex //管理连接读写的锁

}

func (cm *ConnManager) Add(conn siface.IConnection) {

	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", cm.Len())

}

func (cm *ConnManager) Remove(conn siface.IConnection) {

	//保护共享资源Map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除连接信息
	delete(cm.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", cm.Len())
}

//利用ConnID获取链接
func (cm *ConnManager) Get(connID uint32) (siface.IConnection, error) {
	//保护共享资源Map 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

//清除并停止所有连接
func (cm *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range cm.connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", cm.Len())
}

func NewConnManager() siface.IConnManager {

	return &ConnManager{
		connections: make(map[uint32]siface.IConnection),
	}
}
