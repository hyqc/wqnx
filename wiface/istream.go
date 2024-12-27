package wiface

import "github.com/quic-go/quic-go"

type IStream interface {
	GetStream() quic.Stream
	SendMsg(msgId uint32, data []byte) error
	Start()
	Stop()
}
