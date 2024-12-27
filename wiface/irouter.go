package wiface

type IRouter interface {
	Before(request IRequest)
	Handle(request IRequest) (body []byte)
	After(request IRequest)
	GetReqMsgId() uint32
	GetRespMsgId() uint32
}
