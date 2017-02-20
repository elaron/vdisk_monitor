package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"errors"
	"math/rand"
	"vdisk_monitor/src/common/messageIntf"
)

func handleAddVdiskFbMsg(m map[string]string, fbMsg string) error{
		
	opResult, ok := m["OpResult"]
	if false == ok {
		s := fmt.Sprintf("Feedback message lack of OpResult tag! Msg(%s)", fbMsg)
		return errors.New(s)
	}

	msgBody, ok := m["Body"]
	if false == ok {
		s := fmt.Sprintf("Feedback message lack of body! Msg(%s)", fbMsg)
		return errors.New(s)
	}

	if "FAIL" == opResult {
		s := fmt.Sprintf("Add vdisk fail! Err:%s", msgBody)
		return errors.New(s)
	}

	//vdiskId := msgBody
	//startOriginator(vdiskId)
	return nil
}

func feedbackMsgHandler(fbMsg string) error {
	
	var m map[string]string
	
	err := json.Unmarshal([]byte(fbMsg), &m)
	if nil != err {
		s := fmt.Sprintf("Handle feedback msg fail, Err:%s", err.Error())
		return errors.New(s)
	}

	feedBackType, ok := m["MsgType"]
	if false == ok {
		s := fmt.Sprintf("Message lack of MsgType tag, cannot be handled! Msg(%s)", fbMsg)
		return errors.New(s)
	}
	opResult, ok := m["OpResult"]
	if false == ok {
		s := fmt.Sprintf("Message lack of OpResult tag, cannot be handled! Msg(%s)", fbMsg)
		return errors.New(s)
	}

	switch feedBackType {
		case "ADD_VDISK_FEEDBACK":
			//handleAddVdiskFbMsg(m,fbMsg)
			fmt.Println(fbMsg)

		case "REGISTER_AGENT_FEEDBACK":
			
			if opResult == "SUCCESS" {
				runAgent()

			}else {
				fmt.Println(fbMsg)	
			}
			
		default:
			fmt.Println(fbMsg)
	}

	return nil
}

func listenFeedback(conn net.Conn) {
	
	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
	
		//Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))

		feedbackMsgHandler(string(buffer[:n]))

		buffer = make([]byte, 2048)
	}
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

	go listenFeedback(conn)

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
	message := messageIntf.RegisterAgentMsg{
		MsgType: "REGISTER_AGENT",
		Hostname: "aaa",
		Ip:       g_agentConfig.HostIp,
		Id:       g_agentConfig.AgentId,
		PeerAgentId: g_agentConfig.PeerAgentId,
		TcpServerPort: g_agentConfig.TcpServerPort,
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
	registerAgent := messageIntf.RegisterAgentMsg{
		MsgType: "AGENT_HEART_BEAT",
		Hostname: "aaa",
		Ip:       "192.168.56.104",
		Id:       g_agentConfig.AgentId,
		TcpServerPort: g_agentConfig.TcpServerPort,
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
		time.Sleep(3 * time.Second)
	}
}



func sendAddVdiskMsgToMds() {

	sendFunc := sendMsgToMds()
	
	for {

		vmIdx := rand.Intn(30)
		pathIdx := rand.Intn(100)

		vmName := fmt.Sprintf("vm_case%d", vmIdx)
		path := fmt.Sprintf("/root/wyd/case%d/vdisk_%d.qcow2", vmIdx, pathIdx)

		message := messageIntf.AddVdiskMsg{
			MsgType: "ADD_VDISK",
			AgentId: g_agentConfig.AgentId,
			VmId: vmName,
			Path: path,
		}

		b, err := json.Marshal(message)
		if nil != err {
			fmt.Printf("Generate add vdisk msg fail, err :%s", err.Error())
		}

		err = sendFunc(string(b))
		if err != nil {
			fmt.Println("Send addvdisk fail")
			break
		}
		fmt.Println(string(b))
		
		time.Sleep(30 * time.Second)
	}
}
