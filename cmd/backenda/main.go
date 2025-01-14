package main

import (
	"http-reverse-proxy/pkg/server"
)

func main() {

	if err := server.StartServer("configs/backenda.yaml"); err != nil {
		panic(err)
	}
}
