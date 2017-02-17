package main

import (
	"fmt"
	"net"
	"log"
	"os"
)

func setuplistener() {

	addr := fmt.Sprintf("localhost:%d", g_agentConfig.TcpServerPort)

	netListen, err := net.Listen("tcp", addr)
	CheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
	
		Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
		//feedback, _ := msgHandler(buffer[:n])

		buffer = make([]byte, 2048)

		//fmt.Println("Feedback:", feedback)
		
		//go sendFeedback(conn, feedback)
	}
}


func Log(v ...interface{}) {
	log.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
