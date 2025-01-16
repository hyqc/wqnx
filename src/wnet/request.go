package wnet

import (
	"wqnx/src/wiface"
)

type Request struct {
	stream wiface.IStream
	data   wiface.IMessage
}

func NewRequest(stream wiface.IStream, data wiface.IMessage) *Request {
	return &Request{
		stream: stream,
		data:   data,
	}
}

func (r *Request) GetStream() wiface.IStream {
	return r.stream
}

func (r *Request) GetMsgData() []byte {
	return r.data.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.data.GetMsgId()
}
