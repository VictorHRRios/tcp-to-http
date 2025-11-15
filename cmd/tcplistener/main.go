package main

import (
	"fmt"
	"log"
	"net"

	"github.com/VictorHRRios/http/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			request.RequestLine.Method,
			request.RequestLine.RequestTarget,
			request.RequestLine.HttpVersion,
		)
		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Printf("Body:\n%s", string(request.Body))
	}
}
