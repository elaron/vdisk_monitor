package main

import (
	"fmt"
	"bytes"
	"reflect"
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

func compareAgent(srcAgent Agent, agentID int32) (bool, string){
	
	agentTest, err := getAgent(agentID)
	if err != nil {
		return false, "Get agent fail"
	}

	if false == reflect.DeepEqual(agentTest, srcAgent) {
		fmt.Printf("Set agent:%s \nget agent:%s\n", srcAgent, agentTest)
		return false, "Get agent fail, different with seted agent"
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

	rslt, msg := compareAgent(agent, 100)
	if false ==  rslt{
		return rslt, msg
	}
	
	agent.BasicInfo.State = LOSS_CONN

	err = updateAgent(agent)
	if nil != err {
		fmt.Printf("Update agent fail! Err:%s", err.Error())
		return false, "Update agent fail"
	}

	rslt, msg = compareAgent(agent, 100)
	if false ==  rslt{
		return rslt, msg
	}

	return true, ""
}

func case4_GetAgentList() (bool, string) {

	_,err := getAgentList()

	if err != nil {
		return false, "Get agent list fail"
	}

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

	rslt, errMsg = case4_GetAgentList()
	if false == rslt {
		fmt.Println("case4_GetAgentList --- Fail, errMsg: ", errMsg)
	}else{
		fmt.Println("case4_GetAgentList --- Pass")	
	}
	
}