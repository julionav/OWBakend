package server

import (
	"bufio"
	"fmt"
	"io"
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
	buff := make([]byte, 1024)

	for {
		// Read a single byte which contains the message length
		size, err := clientReader.ReadByte()
		if err != nil {
			panic(err)
		}

		// Read the full message, or return an error
		_, err = io.ReadAtLeast(clientReader, buff[:int(size)], 1)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Disconnected: " + client.connection.RemoteAddr().String())
			} else {
				// Seems to come here on disconnect
				fmt.Println("Error with connection " + client.connection.RemoteAddr().String() + err.Error())
			}
		}

		// We chopped the first byte to know the message size.
		// We need to reconstruct the full message appending the first byte.
		fullMessage := append([]byte{size}, buff[:int(size)]...)
		message := string(fullMessage)

		server.Messages <- message
		fmt.Println("Added message to server buffer:", message)

		// Notify back to client to resume execution. Otherwise, client will be blocked waiting for io.
		_, err = client.connection.Write(fullMessage)
		if err != nil {
			fmt.Println("Error writing back to client")
			return
		}

		// Clear buffer
		buff = append(buff[:size], buff[size+1:]...)
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
