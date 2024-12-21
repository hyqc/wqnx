package wiface

import "github.com/quic-go/quic-go"

type IConnection interface {
	SetId(connId uint32) //设置连接ID
	GetId() uint32
	GetConn() quic.Connection //获取连接
	Start()
	Stop()
	SendMsg(msgId uint32, data []byte) error //发送回包消息
}
