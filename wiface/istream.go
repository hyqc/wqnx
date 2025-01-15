package wiface

import (
	"context"
	"github.com/quic-go/quic-go"
)

type IStream interface {
	GetStream() quic.Stream
	SendMsg(msgId uint32, data []byte) error
	Start()
	Stop()
	GetCtx() context.Context
	SetAttr(string, interface{})
	GetAttr(string) (interface{}, bool)
	RemoveAttr(string)
}
