package wnet

import "wqnx/wiface"

type BaseRouter struct {
	ReqMsgId  uint32
	RespMsgId uint32
}

func (b *BaseRouter) Before(request wiface.IRequest) {

}

func (b *BaseRouter) Handle(request wiface.IRequest) (body []byte) {
	return make([]byte, 0)
}

func (b *BaseRouter) After(request wiface.IRequest) {

}

func (b *BaseRouter) GetReqMsgId() uint32 {
	return b.ReqMsgId
}

func (b *BaseRouter) GetRespMsgId() uint32 {
	return b.RespMsgId
}
