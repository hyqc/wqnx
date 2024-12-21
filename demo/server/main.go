package main

import "wqnx/wnet"

func main() {

	wnet.NewServer(
		wnet.WithIP("127.0.0.1"),
		wnet.WithPort(6666),
		wnet.WithHost("localhost"),
	).Run()
}
