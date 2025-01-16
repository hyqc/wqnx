package wnet

import (
	"fmt"
	"wqnx/src/wiface"
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
	defer func() {
		if err := recover(); err != nil {
			SysPrintError(fmt.Sprintf("router handle error:%v; msgId:%v", err, request.GetMsgId()))
		}
	}()
	if call, ok := r.routers[request.GetMsgId()]; ok {
		response := NewResponse(request.GetStream(), call)
		call.Before(request, response)
		call.Handle(request, response)
		call.After(request, response)
	} else {
		SysPrintError("router not found, router id: ", request.GetMsgId())
	}
}
