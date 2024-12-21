package wnet

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
)

type Connection struct {
	ctx            context.Context
	conn           quic.Connection //链接
	id             uint32          //链接ID
	msgChan        chan []byte     //无缓冲读写消息通道
	exitStreamChan chan bool       //退出信号
	isClosed       bool            //通道是否关闭
}

func NewConnection(ctx context.Context, conn quic.Connection, id uint32) *Connection {
	return &Connection{
		ctx:            ctx,
		conn:           conn,
		id:             id,
		exitStreamChan: make(chan bool, 1),
		msgChan:        make(chan []byte),
	}
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
	c.exitStreamChan <- true
}

// SendMsg 发送回包消息
func (c *Connection) SendMsg(stream quic.Stream, msgId uint32, data []byte) error {
	defer func() {
		_ = stream.Close()
	}()
	if c.isClosed {
		return fmt.Errorf("connection is closed")
	}
	dp := NewDataPack()
	body, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		SysPrintError("pack error: ", err.Error())
		return err
	}
	_, err = stream.Write(body)
	//c.msgChan <- body
	return err
}

func (c *Connection) readerStream() {
	SysPrintInfo("reader stream starting...")
	defer c.Stop()
	for {
		// 接收数据流
		Print("99999999999999")
		stream, err := c.conn.AcceptStream(c.ctx)
		if err != nil {
			Print("100000000000")
			SysPrintError("conn accept stream failed, error: ", err)
			return
		}
		Print("10000001111111")
		SysPrintInfo(fmt.Sprintf("accept stream remote addr: %s, stream id: %v ", c.conn.RemoteAddr().String(), stream.StreamID()))

		dp := NewDataPack()
		headerData := make([]byte, dp.GetHeadLen())
		if _, err = stream.Read(headerData); err != nil {
			SysPrintError("serve read error: ", err.Error())
			continue
		}
		Print(string(headerData), "   ", headerData)

		msg, err := dp.Unpack(headerData)
		if err != nil {
			SysPrintError("serve unpack error: ", err.Error())
			continue
		}
		Print(msg.GetDataLen())
		var body = make([]byte, msg.GetDataLen())
		if msg.GetDataLen() > 0 {
			if _, err = stream.Read(body); err != nil {
				SysPrintError("serve read error: ", err.Error())
				continue
			}
		}
		msg.SetData(body)
		SysPrintInfo(fmt.Sprintf("stream id: %v, msgId: %v, dataLen: %v, data: %v ", stream.StreamID(), msg.GetMsgId(), msg.GetDataLen(), string(body)))
		// 将数据发送到消息队列
		go func(stream quic.Stream) {
			err := c.SendMsg(stream, msg.GetMsgId(), msg.GetData())
			if err != nil {
				SysPrintError("send msg error: ", err.Error())
				return
			}
		}(stream)
	}
}
