package main

import (
	"fmt"
	"os"
	"os/signal"
	"wqnx/src/wnet"
)

func main() {
	wnet.MockClient("127.0.0.1", "localhost", 6666)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		fmt.Println("exit")
	}
}
