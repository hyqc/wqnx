package wiface

type IRouter interface {
	Before(request IRequest, response IResponse)
	Handle(request IRequest, response IResponse)
	After(request IRequest, response IResponse)
	GetReqMsgId() uint32
	GetRespMsgId() uint32
}
