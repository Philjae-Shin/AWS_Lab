package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func read(conn net.Conn) {
	//TODO In a continuous loop, read a message from the server and display it.
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Server response:", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		handleError(err)
	}
}

func write(conn net.Conn) {
	//TODO Continually get input from the user and send messages to the server.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter message: ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
		}
		_, err = fmt.Fprintf(conn, msg)
		if err != nil {
			handleError(err)
		}
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	//TODO Try to connect to the server
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		handleError(err)
	}
	defer conn.Close()
	fmt.Println("Client connected")

	//TODO Start asynchronously reading and displaying messages
	go read(conn)

	//TODO Start getting and sending user messages.
	write(conn)
}
