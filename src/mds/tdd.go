package main

import (
	"fmt"
	"bytes"
	"errors"
	"reflect"
)

func case1_EtcdCRUD() error {

	createEtcdValue, deleteEtcdValue, updateEtcdValue, getEtcdValue := createKey(), deleteKey(), updateKey(), getKey()

	uuid := genUUID()
	
	key := bytes.Buffer{}
	key.WriteString(uuid)
	key.WriteString("elar")

	err := createEtcdValue(key.String(), "Hi, nice to see you")
	if err != nil {
		return errors.New("set key fail!")
	
	}else{

		value, err := getEtcdValue(key.String())
		if (err != nil) || (value != "Hi, nice to see you"){
			return errors.New("get key fail")
		}

		err = updateEtcdValue(key.String(), "i am new info")
		if err != nil {
			return errors.New("update key fail")
		}

		value, err = getEtcdValue(key.String())
		if (err != nil) || (value != "i am new info"){
			return errors.New("get key fail")
		}

		err = deleteEtcdValue(key.String())
		if err != nil {
			return errors.New("delete key fail")
		}
	}
	return nil
}

func compareAgent(srcAgent Agent, agentID string) error{
	
	agentTest, err := getAgent(agentID)
	if err != nil {
		return errors.New("Get agent fail")
	}

	if false == reflect.DeepEqual(agentTest, srcAgent) {
		fmt.Printf("Set agent:%s \nget agent:%s\n", srcAgent, agentTest)
		return errors.New("Get agent fail, different with seted agent")
	}

	return nil
}

func case2_AgentCRUD() error {

	var agentID string = "101"
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
		return errors.New("Add agent fail!!!!")
	}

	err = compareAgent(agent, agentID)
	if err !=  nil{
		return err
	}
	
	agent.BasicInfo.State = LOSS_CONN

	err = updateAgent(agent)
	if nil != err {
		fmt.Printf("Update agent fail! Err:%s", err.Error())
		return errors.New("Update agent fail")
	}

	err = compareAgent(agent, agentID)
	if err !=  nil{
		return err
	}

	return nil
}

func case3_addAgent() error {

	agent := Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "agent100",
			Id:         "101",
		},
	}

	deleteAllAgents()

	rslt, errMsg := addAgent(agent)
	
	if rslt != true {
		fmt.Println("Add agent fail, err:", errMsg)
		return errors.New("Add agent fail!")
	}

	agent.BasicInfo.Id = "102"
	rslt, _ = addAgent(agent)
	if rslt != false {
		return errors.New("Detect duplicate AgentIP fail!")
	}

	agent.BasicInfo.Id = "101"
	agent.BasicInfo.HostIp = "10.25.26.47"
	rslt, _ = addAgent(agent)
	if rslt != false {
		return errors.New("Detect duplicate AgentID fail!")
	}

	return nil
}

func case4_removeAgent() error {
	
	agent := Agent{
		BasicInfo: AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "agent100",
			Id:         "101",
		},
	}

	fmt.Println("here 1")
	deleteAllAgents()

	fmt.Println("here 2")
	rslt, errMsg := addAgent(agent)
	
	if rslt != true {
		fmt.Println("Add agent fail, err:", errMsg)
		return errors.New("Add agent fail!")
	}

	fmt.Println("here 3")
	rslt, errMsg = removeAgent("102")
	if false != rslt {
		return errors.New("Fail of detect non-exist agent")
	}

	fmt.Println("here 4")
	rslt, errMsg = removeAgent("101")
	if true != rslt {
		return errors.New("Fail of removing agent")
	}

	return nil
}

func case5_addVdisk() error{
	
	err := addVdisk("101", "vm_case5", "root/case5/os_vdisk.qcow2")
	return err
}

func main() {
	
	var err error

	err = case1_EtcdCRUD()
	if nil != err {
		fmt.Println("case1_EtcdCRUD --- Fail, ", err.Error())
	}else{
		fmt.Println("case1_EtcdCRUD --- Pass")	
	}
	

	err = case2_AgentCRUD()
	if nil != err {
		fmt.Println("case2_AgentCRUD --- Fail, ", err.Error())
	}else{
		fmt.Println("case2_AgentCRUD --- Pass")	
	}

	err = case3_addAgent()
	if nil != err {
		fmt.Println("case3_addAgent --- Fail, ", err.Error())
	}else{
		fmt.Println("case3_addAgent --- Pass")	
	}

	err = case4_removeAgent()
	if nil != err {
		fmt.Println("case4_removeAgent --- Fail, ", err.Error())
	}else{
		fmt.Println("case4_removeAgent --- Pass")	
	}
	
	err = case5_addVdisk()
	if nil != err {
		fmt.Println("case5_addVdisk --- Fail, ", err.Error())
	}else{
		fmt.Println("case5_addVdisk --- Pass")
	}
}