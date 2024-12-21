package wnet

type Message struct {
	id      uint32
	data    []byte
	dataLen uint32
}

func NewDefaultMessage() *Message {
	return &Message{
		data: make([]byte, 0),
	}
}

func NewMessage(msgId uint32, data []byte) *Message {
	return &Message{
		id:      msgId,
		data:    data,
		dataLen: uint32(len(data)),
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetData(data []byte) {
	m.data = data
}
