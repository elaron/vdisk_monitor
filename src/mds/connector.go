package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"encoding/json"
	"errors"
	"vdisk_monitor/src/common/messageIntf"
	"vdisk_monitor/src/common/dbConn"
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

func handleRegisterAgentMsg(m map[string]interface{}, msg string) (feedback string, err error){

	value, ok := m["Hostname"]
	if false == ok {
		s := fmt.Sprintf("Lack of Hostname, register agent fail. Msg:%s", string(msg))
		feedback = genCommonMsgFeedback("REGISTER_AGENT", "FAIL", s)
		err = errors.New(s)
		return
	}
	hostName, ok := value.(string)
	if false == ok {
		s := fmt.Sprintf("Hostname type assertion fail")
		feedback = genCommonMsgFeedback("REGISTER_AGENT", "FAIL", s)
		err = errors.New(s)
		return
	}

	value, ok = m["Ip"]
	if false == ok {
		s := fmt.Sprintf("Lack of Ip, register agent fail. Msg:%s", string(msg))
		feedback = genCommonMsgFeedback("REGISTER_AGENT", "FAIL", s)
		err = errors.New(s)
		return
	}
	hostIP, ok := value.(string)
	if false == ok {
		s := fmt.Sprintf("hostIP type assertion fail")
		feedback = genCommonMsgFeedback("REGISTER_AGENT", "FAIL", s)
		err = errors.New(s)
		return
	}

	value, ok = m["Id"]
	if false == ok {
		s := fmt.Sprintf("Lack of Id, register agent fail. Msg:%s", string(msg))
		feedback = genCommonMsgFeedback("REGISTER_AGENT", "FAIL", s)
		err = errors.New(s)
		return
	}
	agentID, ok := value.(string)
	if false == ok {
		s := fmt.Sprintf("AgentId type assertion fail")
		feedback = genCommonMsgFeedback("REGISTER_AGENT", "FAIL", s)
		err = errors.New(s)
		return
	}


	err = addAgent(agentID, hostIP, hostName)
	if nil != err {
		s := fmt.Sprintf("Register agent fail!")
		feedback = genCommonMsgFeedback("REGISTER_AGENT", "FAIL", s)
		err = errors.New(s)
		return
	}

	fmt.Println("Register agent success")
	
	feedback = genCommonMsgFeedback("REGISTER_AGENT", "SUCCESS", "")
	err = nil

	return
}

func handleAddVdiskMsg(m map[string]interface{}, msg string) (feedback string, err error){
	
	value,ok := m["AgentId"]
	if false == ok {
		s := fmt.Sprintf("Lack of AgentId, add vdisk fail. Msg: ", string(msg))
		feedback = genCommonMsgFeedback("ADD_VDISK", "FAIL", s)
		err = errors.New(s)
		return
	}
	agentID, ok := value.(string)
	if false == ok {
		s := fmt.Sprintf("AgentId type assertion fail")
		feedback = genCommonMsgFeedback("ADD_VDISK", "FAIL", s)
		err = errors.New(s)
		return
	}

	value, ok = m["VmId"]
	if false == ok {
		s := fmt.Sprintf("Lack of VmId, add vdisk fail. Msg:%s", string(msg))
		feedback = genCommonMsgFeedback("ADD_VDISK", "FAIL", s)
		err = errors.New(s)
		return
	}
	vmId, ok := value.(string)
	if false == ok {
		s := fmt.Sprintf("VmId type assertion fail")
		feedback = genCommonMsgFeedback("ADD_VDISK", "FAIL", s)
		err = errors.New(s)
		return
	}

	value, ok = m["Path"]
	if false == ok {
		s := fmt.Sprintf("Lack of Path, add vdisk fail. Msg:%s", string(msg))
		feedback = genCommonMsgFeedback("ADD_VDISK", "FAIL", s)
		err = errors.New(s)
		return
	}
	path, ok := value.(string)
	if false == ok {
		s := fmt.Sprintf("Path type assertion fail")
		feedback = genCommonMsgFeedback("ADD_VDISK", "FAIL", s)
		err = errors.New(s)
		return
	}

	vdiskId, err := addVdisk(agentID, vmId, path)
	if nil != err {
		feedback = genCommonMsgFeedback("ADD_VDISK", "FAIL", err.Error())
		fmt.Println("Add vdisk fail")
		return
	}

	fmt.Println("Add vdisk success")
	feedback = genCommonMsgFeedback("ADD_VDISK", "SUCCESS", vdiskId)
	err = nil
	return
}

func genCommonMsgFeedback(msgType string, opResult string, msgBody string) string{
	
	feedbackMsgType := fmt.Sprintf("%s_FEEDBACK", msgType)
	
	fbMsg := messageIntf.FeedbackMsg{
		MsgType: feedbackMsgType,
		OpResult: opResult,
		Body: msgBody,
	}

	b, convertErr := json.Marshal(fbMsg)
	if nil != convertErr {
		return ""
	}

	return string(b)
}

func msgHandler(jsonMsg []byte) (feedback string, err error){

	var f interface{}

	err = json.Unmarshal(jsonMsg, &f)
	
	if nil != err {
	
		s := fmt.Sprintf("Handle msg(%s) fail, err:", string(jsonMsg), err.Error())	
		err = errors.New(s)
		feedback = genCommonMsgFeedback("UNKNOWN_MSG", "FAIL", s)
		return 
	}

	m := f.(map[string]interface{})

	msgType, ok := m["MsgType"]
	if false == ok {
		s := fmt.Sprintf("Msg dosen't has MsgType tag, msg:%s", string(jsonMsg))
		err = errors.New(s)
		feedback = genCommonMsgFeedback(msgType.(string), "FAIL", s)
		return 
	}

	switch msgType {

		case "REGISTER_AGENT":
			feedback, err = handleRegisterAgentMsg(m, string(jsonMsg))	

		case "ADD_VDISK":
			feedback, err = handleAddVdiskMsg(m, string(jsonMsg))

		case "AGENT_HEART_BEAT":
			fmt.Println(string(jsonMsg))

		default:
			s := fmt.Sprintf("Unrecognize msgType:%s", msgType)
			err = errors.New(s)
			fmt.Println(s)

			feedback = genCommonMsgFeedback("UNKNOWN_MSG", "FAIL", s)

	}
/*
	for key, value := range m {
		fmt.Printf("key:%s value:%s\n", key, value)
	}
*/
	return 
}

func listenFeedback(conn net.Conn) {
	
}

func sendFeedback(conn net.Conn, feedback string) {
	
	conn.Write([]byte(feedback))
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
		feedback, _ := msgHandler(buffer[:n])

		buffer = make([]byte, 2048)

		fmt.Println("Feedback:", feedback)
		
		go sendFeedback(conn, feedback)
	}
}

func connectAgent(agentID string) (net.Conn){

	agent, _ := common.GetAgent(agentID)
	
	server := fmt.Sprintf("%s:%d", agent.BasicInfo.HostIp, agent.BasicInfo.TcpServerPort)

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

func sendMsgToMds(agentID string) func(string) error{
	
	conn := connectAgent(agentID)

	return func(msg string) error {
	
		_,err := conn.Write([]byte(msg))

		if  err != nil {
			fmt.Println(err.Error())
		}
		return err
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
