package wnet

import (
	"bytes"
	"encoding/binary"
	"wqnx/wiface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// id + dataLen 类型都是 uint32
	return 8
}

// Pack 数据包封包
func (d *DataPack) Pack(msg wiface.IMessage) ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	if err := binary.Write(buff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Unpack 数据包解包
func (d *DataPack) Unpack(data []byte) (wiface.IMessage, error) {
	buffReader := bytes.NewReader(data)
	msg := NewDefaultMessage()
	if err := binary.Read(buffReader, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}
	if err := binary.Read(buffReader, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(buffReader, binary.LittleEndian, &msg.data); err != nil {
		return nil, err
	}
	return msg, nil
}
