package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"math/rand"
	//"runtime"
)

type RegisterAgentMsg struct {
	MsgType	 string
	Hostname string
	Ip       string
	Id       string
}

type AddVdiskMsg struct {
	MsgType string
	AgentId string
	VmId string
	Path string
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

func sendRegisteAgentMsg() {
	message := RegisterAgentMsg{
		MsgType: "REGISTER_AGENT",
		Hostname: "aaa",
		Ip:       "192.168.56.104",
		Id:       "100",
	}

	msg, err := json.Marshal(message)
	if nil != err {
		fmt.Println(err.Error())
		return
	}

	sendFunc := sendMsgToMds()
	sendFunc(string(msg))
}

func getAgentIdentifyInfo() string {
	registerAgent := RegisterAgentMsg{
		MsgType: "AGENT_HEART_BEAT",
		Hostname: "aaa",
		Ip:       "192.168.56.104",
		Id:       "100",
	}

	msg, err := json.Marshal(registerAgent)
	if nil != err {
		fmt.Println("encode to json fail!")
	} else {
		fmt.Println("json-body:", string(msg))
	}

	return string(msg)
}

func heartbeatToMds() {

	info := getAgentIdentifyInfo()
	sendHbMsg := sendMsgToMds()
	
	for {
		//runtime.Gosched()
		sendHbMsg(info)
		fmt.Println(info)
		time.Sleep(1 * time.Second)
	}
}



func sendAddVdiskMsgToMds() {

	sendFunc := sendMsgToMds()
	
	for {

		vmIdx := rand.Intn(30)
		pathIdx := rand.Intn(100)

		vmName := fmt.Sprintf("vm_case%d", vmIdx)
		path := fmt.Sprintf("/root/wyd/case%d/vdisk_%d.qcow2", vmIdx, pathIdx)

		message := AddVdiskMsg{
			MsgType: "ADD_VDISK",
			AgentId: "100",
			VmId: vmName,
			Path: path,
		}

		b, err := json.Marshal(message)
		if nil != err {
			fmt.Printf("Generate add vdisk msg fail, err :%s", err.Error())
		}

		sendFunc(string(b))
		time.Sleep(3 * time.Second)
	}
	
}
