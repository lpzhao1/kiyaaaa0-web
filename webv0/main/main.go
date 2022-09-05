package main

import (
	"webv0"
	"webv0/demo"
)

func main() {
	testServer := webv0.NewSdkHttpServer("testServer0")
	testServer.Route("/main", demo.Main)
	testServer.Route("/signup", demo.SignUp)
	testServer.Start(":8080")
}
