package wnet

import (
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"time"
	"wqnx/src/wiface"
)

func MockClient(ip, host string, port int) {
	tlConf := generateTLSConfig(ip, host)
	tlConf.InsecureSkipVerify = true
	ctx := context.Background()
	conn, err := quic.DialAddr(ctx, fmt.Sprintf("%s:%d", ip, port), tlConf, &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		panic(err)
	}

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	var send = 0
	go func() {
		for {
			send++

			dp := NewDataPack()
			body, err := dp.Pack(NewMessage(1, []byte(fmt.Sprintf("你好%v：%v", send, time.Now().Format(time.DateTime)))))
			if err != nil {
				fmt.Println("pack error: ", err)
				return
			}
			if _, err = stream.Write(body); err != nil {
				fmt.Println("client write data error: ", err)
				break
			}
			fmt.Println("client send data: ", string(body))
			msg, err := MockClientReceive(stream)
			if err != nil {
				fmt.Println("client receive data error: ", err)
				goto sleep
			}
			fmt.Println(fmt.Sprintf("client receive data,reqMsgId: %v, dataLen: %v, data: %v ", msg.GetMsgId(), msg.GetDataLen(), string(msg.GetData())))
		sleep:
			{
				time.Sleep(time.Second * 5)
			}
		}
	}()
	select {
	case <-ctx.Done():
		break
	}
}

func MockClientReceive(stream quic.Stream) (wiface.IMessage, error) {
	dp := NewDataPack()
	headerData := make([]byte, dp.GetHeadLen())
	if _, err := stream.Read(headerData); err != nil {
		return nil, err
	}
	msg, err := dp.Unpack(headerData)
	if err != nil {
		return nil, err
	}
	var body = make([]byte, msg.GetDataLen())
	if msg.GetDataLen() > 0 {
		if _, err = stream.Read(body); err != nil {
			return nil, err
		}
	}
	msg.SetData(body)
	return msg, nil
}
