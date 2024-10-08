package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"time"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

var servers = []string{
	"localhost:8030",
	"localhost:8031",
	"localhost:8032",
}

func makeCall(client *rpc.Client, message string, handler string) {
	request := stubs.Request{Message: message}
	response := new(stubs.Response)
	err := client.Call(handler, request, response)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Responded:", response.Message)
}

func main() {
	wordlist := flag.String("file", "wordlist", "File containing words to reverse")
	flag.Parse()

	// wordlist 파일을 열기
	file, err := os.Open(*wordlist)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 서버 리스트의 각 서버에 대해 클라이언트 설정
	clients := []*rpc.Client{}
	for _, server := range servers {
		client, err := rpc.Dial("tcp", server)
		if err != nil {
			fmt.Println("Error connecting to server:", server, "-", err)
			continue
		}
		defer client.Close()
		clients = append(clients, client)
	}

	// 파일에서 각 단어를 읽어 순차적으로 서버에 분배
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		if len(clients) == 0 {
			fmt.Println("No servers available for processing")
			return
		}
		word := scanner.Text()
		client := clients[i%len(clients)] // 서버를 순환하면서 선택
		fmt.Printf("Sending '%s' to server %s\n", word, servers[i%len(servers)])
		makeCall(client, word, stubs.ReverseHandler)
		i++
		time.Sleep(time.Millisecond * 500) // 작업 간격 설정
	}

	// 프리미엄 함수 FastReverse 테스트
	rand.Seed(time.Now().UnixNano())
	if len(clients) > 0 {
		premiumClient := clients[rand.Intn(len(clients))]
		fmt.Println("\n--- Premium Tier ---")
		makeCall(premiumClient, "Premium Service Test", stubs.PremiumReverseHandler)
	}
}
