package wiface

type IRequest interface {
	GetStream() IStream
	GetMsgId() uint32
	GetMsgData() []byte
}
