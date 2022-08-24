package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Server struct {
	sockets  []*Socket
	Messages chan string
}

type Socket struct {
	server               *Server
	connection           net.Conn
	sendChan             chan string
	lastMessageTimestamp time.Time
	disconnectChan       chan bool
}

func (s *Socket) readMessage() (message string, err error) {
	message, err = bufio.NewReader(s.connection).ReadString('\n')
	return
}

func (s *Socket) disconnect() {
	s.disconnectChan <- true
}

func (s *Socket) send(message string) {
	s.sendChan <- message
}

func (s *Socket) IP() string {
	return s.connection.RemoteAddr().String()
}

func (s *Socket) listen() {
	fmt.Printf("Serving %s\n", s.IP())
	newMessageChan := make(chan string)

	// Listen for new messages
	go func() {
		for {
			var message, _ = s.readMessage()
			newMessageChan <- message
		}
	}()

	for {
		select {
		case message := <-newMessageChan:
			s.lastMessageTimestamp = time.Now()
			s.server.Messages <- message
		case message := <-s.sendChan:
			s.connection.Write([]byte(message + "\n"))
		case <-s.disconnectChan:
			s.connection.Close()
			return
		}
	}

}

func NewServer() *Server {
	return &Server{
		Messages: make(chan string),
	}
}

func handleInnactiveSockets(s *Server) {
	fmt.Println("Checking for inactive sockets")
	now := time.Now()
	for _, socket := range s.sockets {
		if socket.lastMessageTimestamp.Add(5 * time.Second).Before(now) {
			// TODO: Lock and remove socket from array remove socket --> use mutex and that weird stuff.
			fmt.Println("Disconnecting inactive socket with ip " + socket.IP())
			socket.disconnect()
		} else {
			fmt.Println("Sending ping to... " + socket.IP())
			socket.send("PING")
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

	// Ping/Pong management
	pingTicker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-pingTicker.C
			handleInnactiveSockets(s)
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		socket := &Socket{server: s, connection: conn, sendChan: make(chan string), lastMessageTimestamp: time.Now()}
		s.sockets = append(s.sockets, socket)
		go socket.listen()
	}
}
