package wnet

import (
	"sync"
	"wqnx/wiface"
)

// ConnectionMgr Connection的管理模块
type ConnectionMgr struct {
	rw          *sync.RWMutex
	connections map[uint32]wiface.IConnection
}

func NewConnectionMgr() *ConnectionMgr {
	return &ConnectionMgr{
		rw:          new(sync.RWMutex),
		connections: make(map[uint32]wiface.IConnection),
	}
}

func (c *ConnectionMgr) Add(conn wiface.IConnection) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.connections[conn.GetId()] = conn
}

func (c *ConnectionMgr) Remove(conn wiface.IConnection) {
	c.rw.Lock()
	defer c.rw.Unlock()
	delete(c.connections, conn.GetId())
}

func (c *ConnectionMgr) Get(connId uint32) (wiface.IConnection, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	}
	return nil, ErrConnectionNotFound
}

func (c *ConnectionMgr) Len() int {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.connections)
}

func (c *ConnectionMgr) Clear() {
	c.rw.Lock()
	defer c.rw.Unlock()
	for _, conn := range c.connections {
		conn.Stop()
	}
	c.connections = make(map[uint32]wiface.IConnection)
}
