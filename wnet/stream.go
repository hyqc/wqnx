package wnet

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"wqnx/wiface"
)

type Stream struct {
	ctx      context.Context
	stream   quic.Stream
	conn     wiface.IConnection
	isClosed bool
}

func NewStream(ctx context.Context, conn wiface.IConnection, stream quic.Stream) *Stream {
	s := &Stream{
		ctx:    ctx,
		stream: stream,
		conn:   conn,
	}

	conn.GetStreamMgr().Add(s)

	return s
}

func (s *Stream) GetStream() quic.Stream {
	return s.stream
}

// SendMsg 发送回包消息
func (s *Stream) SendMsg(msgId uint32, data []byte) error {
	if s.isClosed {
		return fmt.Errorf("connection is closed")
	}
	dp := NewDataPack()
	body, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		SysPrintError("pack error: ", err.Error())
		return err
	}
	_, err = s.stream.Write(body)
	return err
}

func (s *Stream) Handle() {
	dp := NewDataPack()
	headerData := make([]byte, dp.GetHeadLen())
	if _, err := s.stream.Read(headerData); err != nil {
		SysPrintError("serve read error: ", err.Error())
		return
	}

	msg, err := dp.Unpack(headerData)
	if err != nil {
		SysPrintError("serve unpack error: ", err.Error())
		return
	}

	var body = make([]byte, msg.GetDataLen())
	if msg.GetDataLen() > 0 {
		if _, err = s.stream.Read(body); err != nil {
			SysPrintError("serve read error: ", err.Error())
			return
		}
	}
	msg.SetData(body)

	SysPrintInfo(fmt.Sprintf("stream id: %v, msgId: %v, dataLen: %v, data: %v ", s.stream.StreamID(), msg.GetMsgId(), msg.GetDataLen(), string(body)))

	// 找到处理的路由
	go func() {
		err := s.SendMsg(msg.GetMsgId(), msg.GetData())
		if err != nil {
			SysPrintError("send msg error: ", err.Error())
			return
		}
	}()
}
