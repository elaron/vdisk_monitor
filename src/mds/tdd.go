package main

import (
	"fmt"
	"bytes"
)

func case1_findExistedAgent() (bool, string){

	//create agentList
	var agentList []Agent

	agent1  := Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "aaa",
			Id:         100,
			State:      ACTIVE,	
		},
	}

	agent2  := Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "bbb",
			Id:         101,
			State:      ACTIVE,
		},
	}

	agentList = append(agentList, agent1)
	agentList = append(agentList, agent2)

	var agentTest Agent
	var rslt bool

	//test duplicate agentID
	agentTest = Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "192.25.26.46",
			Hostname:   "lalalala",
			Id:         100,
			State:      ACTIVE,
		},
	}

	rslt = isAgentExist(agentList, agentTest)
	if rslt != true {
		return false, "duplicate agentID"
	}

	//test duplicate IP
	agentTest = Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "lulululu",
			Id:         200,
			State:      ACTIVE,
		},
	}

	rslt = isAgentExist(agentList, agentTest)
	if rslt != true {
		return false, "duplicate hostIP"
	}

	return true, ""
}

func case2_EtcdCRUD() (bool, string) {

	createEtcdValue, deleteEtcdValue, updateEtcdValue, getEtcdValue := createKey(), deleteKey(), updateKey(), getKey()

	uuid := genUUID()
	
	key := bytes.Buffer{}
	key.WriteString(uuid)
	key.WriteString("elar")

	err := createEtcdValue(key.String(), "Hi, nice to see you")
	if err != nil {
		return false, "set key fail!"
	
	}else{

		value, err := getEtcdValue(key.String())
		if (err != nil) || (value != "Hi, nice to see you"){
			return false, "get key fail"
		}

		err = updateEtcdValue(key.String(), "i am new info")
		if err != nil {
			return false, "update key fail"
		}

		value, err = getEtcdValue(key.String())
		if (err != nil) || (value != "i am new info"){
			return false, "get key fail"
		}

		err = deleteEtcdValue(key.String())
		if err != nil {
			return false, "delete key fail"
		}
	}
	return true, ""
}

func case3_AgentCRUD() (bool, string) {

	agent := Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "agent100",
			Id:         100,
		},
	}

	deleteAgent(100)

	err := addAgent(agent)
	if err != nil {
		return false, "Add agent fail!!!!"
	}

	agentTest, err := getAgent(100)
	if err != nil {
		return false, "Get agent fail"
	}

	fmt.Println(agentTest)
	return true, ""
}

func main() {
	
	rslt, errMsg := case1_findExistedAgent()
	if false == rslt {
		fmt.Println("case1_findExistedAgent --- Fail, errMsg: ", errMsg)
	}else{
		fmt.Println("case1_findExistedAgent --- Pass")	
	}
	

	rslt, errMsg = case2_EtcdCRUD()
	if false == rslt {
		fmt.Println("case2_EtcdCRUD --- Fail, errMsg: ", errMsg)
	}else{
		fmt.Println("case2_EtcdCRUD --- Pass")	
	}
	

	rslt, errMsg = case3_AgentCRUD()
	if false == rslt {
		fmt.Println("case3_AgentCRUD --- Fail, errMsg: ", errMsg)
	}else{
		fmt.Println("case3_AgentCRUD --- Pass")	
	}
	
}