package server

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct {
	clients  []*Socket
	Messages chan string
}

type Socket struct {
	connection   net.Conn
	sendChan     chan string
	messagesChan chan string
}

func NewServer() *Server {
	return &Server{
		Messages: make(chan string, 10000),
	}
}

func handleConnection(client *Socket) {
	fmt.Printf("Serving %s\n", client.connection.RemoteAddr().String())
	for {
		message, err := bufio.NewReader(client.connection).ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Disconnected: " + client.connection.RemoteAddr().String())
			} else {
				// Seems to come here on disconnect
				fmt.Println("Error with connection " + client.connection.RemoteAddr().String() + err.Error())
			}
			return
		}

		println("Adding message to channel", message)
		// Stuck here
		client.messagesChan <- message
	}
}

func (s *Server) Start(port string) {
	l, err := net.Listen("tcp4", ":"+port)
	defer l.Close()

	if err != nil {
		panic(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := &Socket{connection: c, messagesChan: s.Messages, sendChan: make(chan string)}
		s.clients = append(s.clients, client)
		go handleConnection(client)
	}
}
