package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"encoding/json"
	"errors"
)

func setuplistener() {

	netListen, err := net.Listen("tcp", "localhost:8877")
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

func handleRegisterAgentMsg(m map[string]interface{}, msg string) {

	value, ok := m["Hostname"]
	if false == ok {
		fmt.Printf("Lack of Hostname, register agent fail. Msg:%s", string(msg))
		return
	}
	hostName, ok := value.(string)
	if false == ok {
		fmt.Printf("Hostname type assertion fail")
		return
	}

	value, ok = m["Ip"]
	if false == ok {
		fmt.Printf("Lack of Ip, register agent fail. Msg:%s", string(msg))
		return
	}
	hostIP, ok := value.(string)
	if false == ok {
		fmt.Printf("hostIP type assertion fail")
		return
	}


	value, ok = m["Id"]
	if false == ok {
		fmt.Printf("Lack of Id, register agent fail. Msg:%s", string(msg))
		return
	}
	agentID, ok := value.(string)
	if false == ok {
		fmt.Printf("AgentId type assertion fail")
		return
	}


	err := addAgent(agentID, hostIP, hostName)
	if nil != err {
		fmt.Println("Register agent fail!")
	}

	fmt.Println("Register agent success")
}

func handleAddVdiskMsg(m map[string]interface{}, msg string) {
	
	value,ok := m["AgentId"]
	if false == ok {
		fmt.Println("Lack of AgentId, add vdisk fail. Msg: ", string(msg))
		return
	}
	agentID, ok := value.(string)
	if false == ok {
		fmt.Printf("AgentId type assertion fail")
		return
	}

	value, ok = m["VmId"]
	if false == ok {
		fmt.Printf("Lack of VmId, add vdisk fail. Msg:%s", string(msg))
		return
	}
	vmId, ok := value.(string)
	if false == ok {
		fmt.Printf("VmId type assertion fail")
		return
	}

	value, ok = m["Path"]
	if false == ok {
		fmt.Printf("Lack of Path, add vdisk fail. Msg:%s", string(msg))
		return
	}
	path, ok := value.(string)
	if false == ok {
		fmt.Printf("Path type assertion fail")
		return
	}

	err := addVdisk(agentID, vmId, path)
	if nil != err {
		fmt.Println("Add vdisk fail")
	}

	fmt.Println("Add vdisk success")
}

func msgHandler(jsonMsg []byte) error{

	var f interface{}

	err := json.Unmarshal(jsonMsg, &f)
	if nil != err {
		s := fmt.Sprintf("Handle msg(%s) fail, err:", string(jsonMsg), err.Error())
		return errors.New(s)
	}

	m := f.(map[string]interface{})

	msgType, ok := m["MsgType"]
	if false == ok {
		s := fmt.Sprintf("Msg dosen't has MsgType, msg:%s", string(jsonMsg))
		return errors.New(s)
	}

	switch msgType {

		case "REGISTER_AGENT":
			handleRegisterAgentMsg(m, string(jsonMsg))	

		case "ADD_VDISK":
			handleAddVdiskMsg(m, string(jsonMsg))

		case "AGENT_HEART_BEAT":
			fmt.Println(string(jsonMsg))

		default:
			fmt.Printf("Unrecognize msgType:%s", msgType)
	}
/*
	for key, value := range m {
		fmt.Printf("key:%s value:%s\n", key, value)
	}
*/
	return nil
}

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
	
		//Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))

		msgHandler(buffer[:n])

		buffer = make([]byte, 2048)
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
