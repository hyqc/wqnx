package wnet

import "fmt"

func ExampleBaseRouter_1() {
	base := NewBaseRouter(1, 2)
	fmt.Println(base.GetReqMsgId())

	// Output:
	// 测试doc.go用例
}

func ExampleBaseRouter_GetReqMsgId() {
	base := NewBaseRouter(1, 2)
	fmt.Println(base.GetReqMsgId())

	// Output:
	// 测试doc.go用例
}
