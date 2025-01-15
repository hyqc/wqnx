package wnet

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"sync"
	"wqnx/wiface"
)

type Stream struct {
	ctx       context.Context
	stream    quic.Stream
	conn      wiface.IConnection
	isClosed  bool
	exit      chan bool
	routerMgr wiface.IRouterMgr
	attr      *sync.Map
}

func NewStream(ctx context.Context, conn wiface.IConnection, stream quic.Stream) *Stream {
	s := &Stream{
		ctx:       ctx,
		stream:    stream,
		conn:      conn,
		isClosed:  false,
		exit:      make(chan bool),
		routerMgr: conn.GetServer().GetRouterMgr(),
		attr:      &sync.Map{},
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

func (s *Stream) Start() {
	go s.handle()
	s.conn.GetServer().CallHookAfterStreamStart(s)
	select {
	case <-s.exit:
		SysPrintInfo("stream stopped, stream id: ", s.stream.StreamID())
		return
	}
}

func (s *Stream) handle() {
	defer s.Stop()
	for {
		dp := NewDataPack()
		headerData := make([]byte, dp.GetHeadLen())
		if _, err := s.stream.Read(headerData); err != nil {
			SysPrintWarn("stream read error: ", err.Error())
			return
		}

		msg, err := dp.Unpack(headerData)
		if err != nil {
			SysPrintError("stream unpack error: ", err.Error())
			return
		}

		var body = make([]byte, msg.GetDataLen())
		if msg.GetDataLen() > 0 {
			if _, err = s.stream.Read(body); err != nil {
				SysPrintError("stream read error: ", err.Error())
				return
			}
		}
		msg.SetData(body)

		SysPrintInfo(fmt.Sprintf("stream id: %v, ReqMsgId: %v, dataLen: %v, data: %v ", s.stream.StreamID(), msg.GetMsgId(), msg.GetDataLen(), string(body)))
		go s.routerMgr.Handle(NewRequest(s, msg))
	}
}

func (s *Stream) Stop() {
	if s.isClosed {
		return
	}
	s.isClosed = true

	s.conn.GetServer().CallHookBeforeStreamStop(s)

	s.exit <- true
	s.conn.GetStreamMgr().Remove(s)
}

func (s *Stream) GetCtx() context.Context {
	return s.ctx
}

func (s *Stream) GetAttr(field string) (interface{}, bool) {
	return s.attr.Load(field)
}

func (s *Stream) SetAttr(field string, value interface{}) {
	s.attr.Store(field, value)
}

func (s *Stream) RemoveAttr(field string) {
	s.attr.Delete(field)
}
