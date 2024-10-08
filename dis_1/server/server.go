package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		fmt.Println("Accepted connection")
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited(구분) by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	defer client.Close()
	scanner := bufio.NewScanner(client)
	for scanner.Scan() {
		msg := scanner.Text()
		tidyMsg := strings.TrimSpace(msg)
		msgs <- Message{sender: clientid, message: tidyMsg}
	}
	if err := scanner.Err(); err != nil {
		handleError(err)
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.

	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		handleError(err)
	}
	defer ln.Close()
	fmt.Println("Listening on", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	clientID := 0

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client
			clientID++
			clients[clientID] = conn
			go handleClient(conn, clientID, msgs)
		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			fmt.Printf("Message from client : %s\n", msg.message)
			for id, conn := range clients {
				if id != msg.sender {
					_, err := fmt.Fprintln(conn, msg.sender, msg.message)
					if err != nil {
						handleError(err)
						conn.Close()
						delete(clients, id)
					}
				}
			}
		}
	}
}
