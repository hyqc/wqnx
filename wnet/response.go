package wnet

import "wqnx/wiface"

type Response struct {
	stream wiface.IStream
	router wiface.IRouter
}

func NewResponse(stream wiface.IStream, router wiface.IRouter) *Response {
	return &Response{
		stream: stream,
		router: router,
	}
}

func (r *Response) GetStream() wiface.IStream {
	return r.stream
}

func (r *Response) SendMsg(body []byte) error {
	return r.stream.SendMsg(r.router.GetRespMsgId(), body)
}
