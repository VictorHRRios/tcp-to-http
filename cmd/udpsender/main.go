package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address := "localhost:42069"
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Print("an error has occurred", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Print("an error has occurred", err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Print("> ")
		read, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(read))
		if err != nil {
			fmt.Print(err)
		}
	}
}
