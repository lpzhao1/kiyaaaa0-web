package main

import (
	"webv1"
	"webv1/demo"
)

func main() {

	testServer := webv1.NewSdkHttpServer("testServer0")
	testServer.Route("/main", demo.Main)
	testServer.Route("/signup", demo.SignUp)
	testServer.Start(":8080")

}
