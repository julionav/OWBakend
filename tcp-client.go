package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	inputChan := make(chan string)
	serverChan := make(chan string)

	inputReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	go func() {
		for {
			text, _ := inputReader.ReadString('\n')
			inputChan <- text
		}
	}()

	go func() {
		for {
			response, _, _ := serverReader.ReadLine()
			serverChan <- string(response)
		}
	}()

	for {
		select {
		case serverMessage := <-serverChan:
			if serverMessage == "PING" {
				fmt.Fprintf(conn, "PONG\n")
			} else {
				fmt.Println("Server response ->: " + serverMessage)
			}
		case inputMessage := <-inputChan:
			fmt.Fprintf(conn, inputMessage)
		}
	}
}
