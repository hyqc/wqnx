package wnet

import "wqnx/wiface"

type BaseRouter struct {
	reqMsgId  uint32
	respMsgId uint32
}

func NewBaseRouter(reqMsgId, respMsgId uint32) wiface.IRouter {
	return &BaseRouter{
		reqMsgId:  reqMsgId,
		respMsgId: respMsgId,
	}
}

func (b *BaseRouter) Before(request wiface.IRequest, response wiface.IResponse) {

}

func (b *BaseRouter) Handle(request wiface.IRequest, response wiface.IResponse) {

}

func (b *BaseRouter) After(request wiface.IRequest, response wiface.IResponse) {

}

func (b *BaseRouter) GetReqMsgId() uint32 {
	return b.reqMsgId
}

func (b *BaseRouter) GetRespMsgId() uint32 {
	return b.respMsgId
}
