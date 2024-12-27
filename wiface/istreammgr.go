package wiface

import "github.com/quic-go/quic-go"

type IStreamMgr interface {
	Add(stream IStream)
	Remove(conn IStream)
	Get(id quic.StreamID) (IStream, error)
	Len() int
	Clear()
}
