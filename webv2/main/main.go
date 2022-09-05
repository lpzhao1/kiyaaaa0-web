package main

import (
	"webv2"
	"webv2/demo"
)

func main() {
	testServer := webv2.NewSdkHttpServer("testServer0")

	testServer.Route("GET", "/main", demo.Main)
	testServer.Route("POST", "/signup", demo.SignUp)

	if err := testServer.Start(":8080"); err != nil {
		panic(err)
	}
}
