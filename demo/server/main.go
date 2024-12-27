package main

import (
	"fmt"
	"time"
	"wqnx/wiface"
	"wqnx/wnet"
)

func main() {

	wnet.NewServer(
		wnet.WithIP("127.0.0.1"),
		wnet.WithPort(6666),
		wnet.WithHost("localhost"),
	).AddRouters(routers...).Run()
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

func (u *UserInfo) Handle(request wiface.IRequest) []byte {
	respBody := []byte(time.Now().Format(time.RFC3339))
	fmt.Println("================ ", time.Now().Format(time.RFC3339))
	return respBody
}
