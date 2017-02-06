package main

import (
	"fmt"
	"net"
	"os"
)

func sender(conn net.Conn) {
	words := "hello world!"
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		conn.Write([]byte(words))
	}
	fmt.Println("send over")

}

func main() {
	server := "127.0.0.1:8877"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	sender(conn)

}
