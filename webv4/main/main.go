package main

import (
	"context"
	"fmt"
	"time"
	"webv4"
	"webv4/demo"
)

func main() {

	shutdown := webv4.NewGracefulShutdown()

	server := webv4.NewSdkHttpServer("my-test-server",
		// 添加两个Filter
		webv4.MetricFilterBuilder, shutdown.ShutdownFilterBuilder)

	adminServer := webv4.NewSdkHttpServer("admin-test-server",
		webv4.MetricFilterBuilder, shutdown.ShutdownFilterBuilder)

	server.Route("POST", "/signup", demo.SignUp)
	server.Route("GET", "/main", demo.Main)

	go func() {
		if err := adminServer.Start(":8081"); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := server.Start(":8080"); err != nil {
			// 服务器都没能成功启动，快速失败
			panic(err)
		}
	}()

	// 先执行 RejectNewRequestAndWaiting，等待所有的请求
	// 然后我们关闭 server，如果是多个 server，可以多个 goroutine 一起关闭
	//
	webv4.WaitForShutdown(
		func(ctx context.Context) error {
			// 假设这里有一个 hook
			// 可以通知网关我们要下线了
			fmt.Println("mock notify gateway")
			time.Sleep(time.Second * 2)
			return nil
		},
		shutdown.RejectNewRequestAndWaiting,
		// 全部请求处理完了我们就可以关闭 server了
		webv4.BuildCloseServerHook(server, adminServer),
		func(ctx context.Context) error {
			// 假设这里要清理一些执行过程中生成的临时资源
			fmt.Println("mock release resources")
			time.Sleep(time.Second * 2)
			return nil
		})

}
