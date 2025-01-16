package wiface

type IMessage interface {
	// GetMsgId 获取消息ID
	GetMsgId() uint32
	// SetMsgId 设置消息ID
	SetMsgId(uint32)
	// GetDataLen 获取消息长度
	GetDataLen() uint32

	// GetData 获取消息内容
	GetData() []byte
	// SetData 设置消息内容
	SetData([]byte)
}
