package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"runtime"
)

type RegisterAgentMsg struct {
	Hostname string
	Ip       string
	Id       int32
}

func connectMds() (net.Conn){

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
	return conn
}

func sendMsgToMds() func(string) error{
	
	conn := connectMds()

	return func(msg string) error {
	
		_,err := conn.Write([]byte(msg))

		if  err != nil {
			fmt.Println(err.Error())
		}
		return err
	}
}

func getAgentIdentifyInfo() string {
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

	return string(msg)
}

func heartbeatToMds(msg string) {

	//info := getAgentIdentifyInfo()
	sendHbMsg := sendMsgToMds()
	
	for {
		runtime.Gosched()
		sendHbMsg(msg)
		fmt.Println(msg)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	info := getAgentIdentifyInfo()
	go heartbeatToMds(info)
	heartbeatToMds("world")
	//heartbeatToMds()
	//heartbeatToMds("{'Hostname':'bbb,'Ip':'192.168.56.104','Id':100}")
}
