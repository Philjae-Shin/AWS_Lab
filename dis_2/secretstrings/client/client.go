package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"sync"
	"time"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

var servers = []string{
	"35.175.153.208:8030",
	"54.81.38.201:8030",
	//"localhost:8032",
}

func makeCall(client *rpc.Client, message string, handler string, wg *sync.WaitGroup) {
	defer wg.Done()

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

	if len(clients) == 0 {
		fmt.Println("No servers available for processing")
		return
	}

	// 파일에서 각 단어를 읽어 순차적으로 서버에 분배
	// Optional: 2 instances working together
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		word := scanner.Text()
		client := clients[i%len(clients)]
		wg.Add(1)
		go makeCall(client, word, stubs.ReverseHandler, &wg)
		i++
	}

	wg.Wait()

	// 프리미엄 함수 FastReverse 테스트
	rand.Seed(time.Now().UnixNano())
	if len(clients) > 0 {
		premiumClient := clients[rand.Intn(len(clients))]
		fmt.Println("\n--- Premium Tier ---")
		wg.Add(1)
		go makeCall(premiumClient, "Premium Service Test", stubs.PremiumReverseHandler, &wg)
	}
	wg.Wait()
	fmt.Println("\n--- All requests complete ---")
}
