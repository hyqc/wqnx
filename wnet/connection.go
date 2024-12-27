package wnet

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"wqnx/wiface"
)

type Connection struct {
	ctx            context.Context
	conn           quic.Connection //链接
	id             uint32          //链接ID
	msgChan        chan []byte     //无缓冲读写消息通道
	exitStreamChan chan bool       //退出信号
	isClosed       bool            //通道是否关闭
	server         wiface.IServer  //链接所属的服务
	streamMgr      wiface.IStreamMgr
}

func NewConnection(ctx context.Context, server wiface.IServer, conn quic.Connection, id uint32) *Connection {
	c := &Connection{
		ctx:            ctx,
		conn:           conn,
		id:             id,
		exitStreamChan: make(chan bool, 1),
		msgChan:        make(chan []byte),
		server:         server,
		streamMgr:      NewStreamMgr(),
	}

	server.GetConnMgr().Add(c)

	return c
}

func (c *Connection) GetId() uint32 {
	return c.id
}

func (c *Connection) SetId(connId uint32) {
	c.id = connId
}

func (c *Connection) GetConn() quic.Connection {
	return c.conn
}

func (c *Connection) Start() {
	SysPrintInfo("connection starting...")
	go c.readerStream()
	//go c.writerStream()
	select {
	case <-c.exitStreamChan:
		SysPrintInfo("connection stopped, connection id: ", c.id)
		return
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.exitStreamChan <- true

	c.server.GetConnMgr().Remove(c)
}

func (c *Connection) readerStream() {
	SysPrintInfo("reader stream starting...")
	defer c.Stop()
	for {
		// 接收数据流
		acceptStream, err := c.conn.AcceptStream(c.ctx)
		if err != nil {
			SysPrintError("conn accept stream failed, error: ", err)
			return
		}
		SysPrintInfo(fmt.Sprintf("accept stream remote addr: %s, stream id: %v ", c.conn.RemoteAddr().String(), acceptStream.StreamID()))

		NewStream(c.ctx, c, acceptStream).Handle()
	}
}

func (c *Connection) GetStreamMgr() wiface.IStreamMgr {
	return c.streamMgr
}
