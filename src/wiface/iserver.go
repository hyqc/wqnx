package wiface

import "context"

type IServer interface {
	Start()                     //开启服务
	Stop()                      //停止服务
	Run()                       //运行服务
	GetConnMgr() IConnectionMgr //获取连接管理
	GetCtx() context.Context
	AddRouters(...IRouter) IServer
	GetRouterMgr() IRouterMgr
	CallHookAfterConnStart(IConnection)
	CallHookAfterStreamStart(IStream)
	CallHookBeforeConnStop(IConnection)
	CallHookBeforeStreamStop(IStream)
}
