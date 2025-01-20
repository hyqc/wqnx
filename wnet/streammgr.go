package wnet

import (
	"github.com/quic-go/quic-go"
	"sync"
	"wqnx/wiface"
)

type StreamMgr struct {
	rw      *sync.RWMutex
	streams map[quic.StreamID]wiface.IStream
}

func NewStreamMgr() *StreamMgr {
	return &StreamMgr{
		rw:      &sync.RWMutex{},
		streams: make(map[quic.StreamID]wiface.IStream),
	}
}

func (s *StreamMgr) Add(stream wiface.IStream) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.streams[stream.GetStream().StreamID()] = stream
}

func (s *StreamMgr) Remove(stream wiface.IStream) {
	s.rw.Lock()
	defer s.rw.Unlock()
	delete(s.streams, stream.GetStream().StreamID())
}

func (s *StreamMgr) Get(id quic.StreamID) (wiface.IStream, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if val, ok := s.streams[id]; ok {
		return val, nil
	}
	return nil, ErrStreamNotFound
}

func (s *StreamMgr) Len() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.streams)
}

func (s *StreamMgr) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	for _, stream := range s.streams {
		_ = stream.GetStream().Close()
	}
	s.streams = make(map[quic.StreamID]wiface.IStream)
}
