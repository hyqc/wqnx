package wiface

type IRouterMgr interface {
	Add(r IRouter)
	Handle(request IRequest)
}
