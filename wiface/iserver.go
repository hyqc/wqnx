package wiface

type IServer interface {
	Start() //开启服务
	Stop()  //停止服务
	Run()   //运行服务
}
