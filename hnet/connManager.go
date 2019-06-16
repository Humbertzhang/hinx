package hnet

import (
	"errors"
	"fmt"
	"hinx/hiface"
	"sync"
)

type ConnManager struct {
	// connection 的集合
	connections map[uint32] hiface.IConnection
	// connections 的锁,读写锁, 用于保护连接集合
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32] hiface.IConnection),
		connLock:    sync.RWMutex{},
	}
}


// 添加连接
func (cm *ConnManager) AddConn(conn hiface.IConnection) {
	// 保护资源map, 上写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将connection加入到 connections中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection ", conn.GetConnID(), " [add] to ConnManager successfully. conn num = ", cm.Len())
}

// 删除
func (cm *ConnManager) DeleteConn(conn hiface.IConnection) {
	// 保护资源map, 上写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection ", conn.GetConnID(), " [deleted] successfully. conn num = ", cm.Len())
}

// 根据ID获取
func (cm *ConnManager) GetConnByID(connID uint32) (hiface.IConnection, error) {
	// 读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	conn, ok := cm.connections[connID]
	if !ok {
		return nil, errors.New("can not find connection:" + string(connID))
	}
	return conn, nil
}

// 获取连接个数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清除并终止所有连接
func (cm *ConnManager) ClearConn() {
	// 保护资源map, 上写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for id, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, id)
	}
	fmt.Println("clear all connections. conn num = ", cm.Len())
}