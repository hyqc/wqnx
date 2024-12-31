package rate

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"testing"
)

func TestNew(t *testing.T) {
	rate := New(context.Background(),
		WithCapacity(1),
		WithRate(1, 5),
	)

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)

	go func() {
		for {
			if rate.Allow() {
				fmt.Println("allow")
			} else {
				fmt.Println("deny")
			}

			if err := rate.WaitN(rate.ctx, 1); err != nil {
				
			}
		}
	}()

	select {
	case <-ch:
		return
	}
}
