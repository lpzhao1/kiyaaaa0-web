package webv4

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

func WaitForShutdown(hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	// Notify 让 signal 包将输入信号转发到 signals
	signal.Notify(signals, ShutdownSignals...)
	select {
	case sig := <-signals:
		fmt.Printf("get signal %s, application will shutdown \n", sig)
		// 十分钟处理时间，结束直接强行退出
		time.AfterFunc(time.Minute*10, func() {
			fmt.Println("Shutdown gracefully timeout, application will shutdown immediately")
			os.Exit(1)
		})

		// 最初版：直接退出
		// os.Exit(0)
		// 成功退出

		// 改造 WaitForShutdown 使其能接收 hook
		for _, h := range hooks {
			ctx, cancel := context.WithTimeout(context.Background(),
				time.Second*30)
			err := h(ctx)
			if err != nil {
				fmt.Printf("failed to run hook, err: %v \n", err)
			}
			cancel()
		}
		os.Exit(0)
	}
}

type GracefulShutdown struct {
	// 还在处理中的请求数
	reqCnt int64
	// closing 大于 1 就说明要关闭了
	closing int32

	// 用channel来通知已经处理完了所有请求
	zeroReqCnt chan struct{}
}

var ErrorHookTimeout = errors.New("the hook timeout")

// RejectNewRequestAndWaiting 将会拒绝新的请求，并且等待处理中的请求
func (g *GracefulShutdown) RejectNewRequestAndWaiting(ctx context.Context) error {
	//
	atomic.AddInt32(&g.closing, 1)

	// 特殊case 关闭之前其实就已经处理完了请求
	if atomic.LoadInt64(&g.reqCnt) == 0 {
		return nil
	}

	done := ctx.Done()
	select {
	case <-done:
		fmt.Println("超时了还没等到所有请求执行完毕")
		return ErrorHookTimeout
	case <-g.zeroReqCnt:
		fmt.Println("全部请求处理完了")
	}
	return nil
}

func (g *GracefulShutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		// 开始拒绝所有的请求
		cl := atomic.LoadInt32(&g.closing)
		if cl > 0 {
			c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		atomic.AddInt64(&g.reqCnt, 1)
		next(c)
		n := atomic.AddInt64(&g.reqCnt, -1)
		// 已经开始关闭了，而且请求数为0，
		if cl > 0 && n == 0 {
			g.zeroReqCnt <- struct{}{}
		}
	}
}

func NewGracefulShutdown() *GracefulShutdown {
	return &GracefulShutdown{
		zeroReqCnt: make(chan struct{}),
	}
}
