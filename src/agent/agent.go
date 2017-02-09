package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type RegisterAgentMsg struct {
	Hostname string
	Ip       string
	Id       int32
}

func heartbeatToMds(conn net.Conn) {

	registerAgent := RegisterAgentMsg{
		Hostname: "aaa",
		Ip:       "192.168.56.104",
		Id:       100,
	}

	msg, err := json.Marshal(registerAgent)
	if nil != err {
		fmt.Println("encode to json fail!")
	} else {
		fmt.Println("json-body:", string(msg))
	}

	for {
		conn.Write([]byte(msg))
		fmt.Println("send over")
		time.Sleep(5 * time.Second)
	}
}
func sayHi(msg string) {
	for i := 0; i < 10; i++ {
		fmt.Println(msg)
	}
}
func main() {
	go sayHi("world")
	sayHi("hello")
}

func connMds() {

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
	go heartbeatToMds(conn)
	fmt.Println("goroutine end")

}
