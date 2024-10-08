// 파일: server.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func ReverseString(s string, delay int) string {
	time.Sleep(time.Duration(rand.Intn(delay)) * time.Second)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type SecretStringOperations struct{}

func (s *SecretStringOperations) Reverse(req stubs.Request, res *stubs.Response) error {
	if req.Message == "" {
		return errors.New("A message must be specified")
	}
	fmt.Println("Received Message:", req.Message)
	res.Message = ReverseString(req.Message, 10)
	return nil
}

func (s *SecretStringOperations) FastReverse(req stubs.Request, res *stubs.Response) error {
	if req.Message == "" {
		return errors.New("A message must be specified")
	}
	res.Message = ReverseString(req.Message, 2)
	return nil
}

func main() {
	port := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	rpc.Register(&SecretStringOperations{})

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port", *port)

	rpc.Accept(listener)
}
