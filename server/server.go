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
		Messages: make(chan string),
	}
}

func handleConnection(server *Server, client *Socket) {
	fmt.Printf("Serving %s\n", client.connection.RemoteAddr().String())

	clientReader := bufio.NewReader(client.connection)

	for {
		message, err := clientReader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Disconnected: " + client.connection.RemoteAddr().String())
			} else {
				// Seems to come here on disconnect
				fmt.Println("Error with connection " + client.connection.RemoteAddr().String() + err.Error())
			}
		}

		server.Messages <- message
		fmt.Println("Added message to server buffer:", message)

		// Notify back to client to resume execution. Otherwise, client will be blocked waiting for io.
		_, err = client.connection.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing back to client")
			return
		}
	}
}

func (s *Server) Start(port string) {
	l, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Println("Server started")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := &Socket{connection: conn, messagesChan: make(chan string), sendChan: make(chan string)}
		s.clients = append(s.clients, client)
		go handleConnection(s, client)
	}
}
