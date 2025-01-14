package main

import (
	"http-reverse-proxy/pkg/server"
)

func main() {

	if err := server.StartServer("configs/backendb.yaml"); err != nil {
		panic(err)
	}
}
