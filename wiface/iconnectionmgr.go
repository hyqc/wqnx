package wiface

// IConnectionMgr 连接管理接口
type IConnectionMgr interface {
	Add(conn IConnection)                   //添加
	Remove(conn IConnection)                //移除
	Get(connId uint32) (IConnection, error) //获取
	Len() int                               //长度
	Clear()                                 //清空
}
