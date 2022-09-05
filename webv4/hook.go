package webv4

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Hook func(ctx context.Context) error

func BuildCloseServerHook(servers ...Server) Hook {
	return func(ctx context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{})
		wg.Add(len(servers))

		for _, s := range servers {
			go func(svr Server) {
				err := svr.Shutdown(ctx)
				if err != nil {
					fmt.Printf("server shutdown error: %v \n", err)
				}
				time.Sleep(time.Second)
				wg.Done()
			}(s)
		}
		go func() {
			wg.Wait()
			doneCh <- struct{}{}
		}()
		select {
		case <-ctx.Done():
			fmt.Printf("closing servers timeout \n")
			return ErrorHookTimeout
		case <-doneCh:
			fmt.Printf("close all servers \n")
			return nil
		}
	}
}
