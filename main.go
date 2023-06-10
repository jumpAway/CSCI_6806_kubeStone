package main

import (
	"kubeStone/m/v2/pkg/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
