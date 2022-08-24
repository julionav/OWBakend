package main

import (
	"fmt"
	"game-service/server"
)

func main() {
	s := server.NewServer()
	go s.Start("3000")
	
	fmt.Println("Server started")

	for {
		m := <-s.Messages
		println("Message from server received! ---> ", m)
	}
}
