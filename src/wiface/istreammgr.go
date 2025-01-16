package wiface

import "github.com/quic-go/quic-go"

type IStreamMgr interface {
	Add(stream IStream)
	Remove(stream IStream)
	Get(id quic.StreamID) (IStream, error)
	Len() int
	Clear()
}
