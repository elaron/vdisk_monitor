package main

import (
	"fmt"
	"bytes"
	"reflect"
)

func case1_EtcdCRUD() (bool, string) {

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

func case2_AgentCRUD() (bool, string) {

	var agentID int32 = 101
	agent := Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "agent100",
			Id:         agentID,
		},
	}

	deleteAgent(agentID)

	err := createAgent(agent)
	if err != nil {
		return false, "Add agent fail!!!!"
	}

	rslt, msg := compareAgent(agent, agentID)
	if false ==  rslt{
		return rslt, msg
	}
	
	agent.BasicInfo.State = LOSS_CONN

	err = updateAgent(agent)
	if nil != err {
		fmt.Printf("Update agent fail! Err:%s", err.Error())
		return false, "Update agent fail"
	}

	rslt, msg = compareAgent(agent, agentID)
	if false ==  rslt{
		return rslt, msg
	}

	return true, ""
}

func case3_addAgent() (bool, string) {

	var agentID int32 = 101
	agent := Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "agent100",
			Id:         agentID,
		},
	}

	deleteAgent(101)
	deleteAgent(100)

	rslt, errMsg := addAgent(agent)
	
	if rslt != true {
		fmt.Println("Add agent fail, err:", errMsg)
		return false, "Add agent fail!"
	}

	agent.BasicInfo.Id = 102
	rslt, _ = addAgent(agent)
	if rslt != false {
		return false, "Detect duplicate AgentIP fail!"
	}

	agent.BasicInfo.Id = 101
	agent.BasicInfo.HostIp = "10.25.26.47"
	rslt, _ = addAgent(agent)
	if rslt != false {
		return false, "Detect duplicate AgentID fail!"
	}

	return true, ""
}

func main() {
	
	rslt, errMsg := case1_EtcdCRUD()
	if false == rslt {
		fmt.Println("case1_EtcdCRUD --- Fail, errMsg: ", errMsg)
	}else{
		fmt.Println("case1_EtcdCRUD --- Pass")	
	}
	

	rslt, errMsg = case2_AgentCRUD()
	if false == rslt {
		fmt.Println("case2_AgentCRUD --- Fail, errMsg: ", errMsg)
	}else{
		fmt.Println("case2_AgentCRUD --- Pass")	
	}

	rslt, errMsg = case3_addAgent()
	if false == rslt {
		fmt.Println("case3_GetAgentList --- Fail, errMsg: ", errMsg)
	}else{
		fmt.Println("case3_GetAgentList --- Pass")	
	}
	
}