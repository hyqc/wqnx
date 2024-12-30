package main

import (
	"fmt"
	"time"
	"wqnx/wiface"
	"wqnx/wnet"
)

func main() {
	wnet.NewServer().AddRouters(routers...).Run()
}

var (
	routers = []wiface.IRouter{
		NewUserInfo(1, 2),
	}
)

type UserInfo struct {
	*wnet.BaseRouter
}

func NewUserInfo(reqMsgId, respMsgId uint32) *UserInfo {
	return &UserInfo{
		&wnet.BaseRouter{
			ReqMsgId:  reqMsgId,
			RespMsgId: respMsgId,
		},
	}
}

func (u *UserInfo) Handle(request wiface.IRequest, response wiface.IResponse) {
	data := []byte(fmt.Sprintf("DateTime: %s, Request MsgId: %v, MsgData: %v", time.Now().Format(time.RFC3339),
		request.GetMsgId(), string(request.GetMsgData())))
	if err := response.SendMsg(data); err != nil {
		fmt.Println("send error: ", err)
	}
}
