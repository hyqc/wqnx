package wnet

import (
	"fmt"
	"wqnx/wiface"
)

type RouterMgr struct {
	routers map[uint32]wiface.IRouter
}

func NewRouterMgr() *RouterMgr {
	return &RouterMgr{
		routers: make(map[uint32]wiface.IRouter),
	}
}

func (r *RouterMgr) Add(i wiface.IRouter) {
	if _, ok := r.routers[i.GetReqMsgId()]; ok {
		panic(fmt.Sprintf("repeated router id: %d", i.GetReqMsgId()))
	}
	r.routers[i.GetReqMsgId()] = i
}

func (r *RouterMgr) Handle(request wiface.IRequest) {
	if call, ok := r.routers[request.GetMsgId()]; ok {
		call.Before(request)
		body := call.Handle(request)
		call.After(request)
		err := NewResponse(request.GetStream(), call).SendMsg(body)
		if err != nil {
			SysPrintError("send msg error: ", err)
		}
	} else {
		SysPrintError("router not found, router id: ", request.GetMsgId())
	}
}
