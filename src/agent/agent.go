package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type RegisterAgentMsg struct {
	Hostname string
	Ip       string
	Id       int32
}

func sender(conn net.Conn, msg string) {
	conn.Write([]byte(msg))
	fmt.Println("send over")

}

func main() {
	registerAgent := RegisterAgentMsg{
		Hostname: "aaa",
		Ip:       "192.168.56.104",
		Id:       100,
	}

	b, err := json.Marshal(registerAgent)
	if nil != err {
		fmt.Println("encode to json fail!")
	} else {
		fmt.Println("json-body:", string(b))
	}

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
	sender(conn, string(b))

}
