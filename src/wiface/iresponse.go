package wiface

type IResponse interface {
	GetStream() IStream
	SendMsg(data []byte) error
}
