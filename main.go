package main

import (
	"game-service/server"
)

func main() {
	s := server.NewServer()
	go s.Start("3000")

	for {
		select {
		case m := <-s.Messages:
			println("Message from server received! ---> ", m)
		}
	}
}
