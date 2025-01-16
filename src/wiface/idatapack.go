package wiface

type IDataPack interface {
	// GetHeadLen 获取包头长度方法
	GetHeadLen() uint32
	// Pack 封包方法(压缩数据)
	Pack(msg IMessage) ([]byte, error)
	// Unpack 解包
	Unpack([]byte) (IMessage, error)
}
