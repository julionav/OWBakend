package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print(">> ")
		text, _ := inputReader.ReadString('\n')
		fmt.Fprintf(conn, text)

		response, _, err := serverReader.ReadLine()
		if err != nil {
			fmt.Println("Error reading line " + err.Error())
		}

		fmt.Println("Server response ->: " + string(response))

		if strings.TrimSpace(text) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
