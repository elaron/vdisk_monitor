package main

import (
	"fmt"
	"time"
	"bytes"
	"errors"
	"reflect"
	"vdisk_monitor/src/common/etcdInterface"
	"vdisk_monitor/src/common/dbConn"
)

func case1_EtcdCRUD() error {

	createEtcdValue, deleteEtcdValue, updateEtcdValue, getEtcdValue := etcdIntf.CreateKey(), etcdIntf.DeleteKey(), etcdIntf.UpdateKey(), etcdIntf.GetKey()

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

func compareAgent(srcAgent common.Agent, agentID string) error{
	
	agentTest, err := common.GetAgent(agentID)
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
	agent := common.Agent{
		BasicInfo: common.AgentBasicInfo {
			HostIp:     "10.25.26.46",
			Hostname:   "agent100",
			Id:         agentID,
		},
	}

	common.DeleteAgent(agentID)

	err := common.CreateAgent(agent)
	if err != nil {
		return errors.New("Add agent fail!!!!")
	}

	err = compareAgent(agent, agentID)
	if err !=  nil{
		return err
	}
	
	agent.BasicInfo.State = common.LOSS_CONN

	err = common.UpdateAgent(agent)
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

	agentID := "101"
	hostIp := "10.25.26.46"
	hostname := "agent100"

	common.DeleteAllAgents()

	err := addAgent(agentID, hostIp, hostname)
	
	if nil != err {
		return err
	}

	agentID = "102"
	err = addAgent(agentID, hostIp, hostname)
	if nil == err {
		return errors.New("Detect duplicate AgentIP fail!")
	}

	agentID = "101"
	hostIp = "10.25.26.47"
	err = addAgent(agentID, hostIp, hostname)
	if err == nil {
		return errors.New("Detect duplicate AgentID fail!")
	}

	return nil
}

func case4_removeAgent() error {
	
	agentID := "101"
	hostIp := "10.25.26.46"
	hostname := "agent100"

	common.DeleteAllAgents()

	err := addAgent(agentID, hostIp, hostname)
	if nil != err {
		return err
	}

	err = removeAgent("102")
	if nil == err {
		return errors.New("Fail of detect non-exist agent")
	}

	err = removeAgent("101")
	if nil != err {
		return errors.New("Fail of removing agent")
	}

	return nil
}

func case5_addVdisk() error{
	
	agentID := "101"
	hostIp := "10.25.26.46"
	hostname := "agent100"

	common.DeleteAllAgents()
	common.DeleteAllVdisks()

	err := addAgent(agentID, hostIp, hostname)
	if nil != err {
		return err
	}

	_, err = addVdisk("101", "vm_case5", "root/case5/os_vdisk.qcow2")
	if nil != err {
		return err
	}

	_, err = addVdisk("101", "vm_case5", "root/case5/os_vdisk.qcow2")
	if nil == err {
		s := fmt.Sprintf("Fail of detecting duplicate vdisk!")
		return errors.New(s)
	}

	t1 := time.Now()

	var loopCnt int = 10
	for i := 0; i < loopCnt; i++ {
		path := fmt.Sprintf("root/case5/vdisk%d.qcow2", i)

		_, err := addVdisk("101", "vm_case5", path)
		if nil != err {
			return err
		}
	}
	
	fmt.Printf("Add %d vdisks need:%q\n ", loopCnt, time.Since(t1))

	t2 := time.Now()

	for i := 0; i < loopCnt; i++ {
		path := fmt.Sprintf("root/case5/vdisk%d.qcow2", i)

		err := removeVdisk("", "101", path)
		if nil != err {
			return err
		}
	}

	fmt.Printf("Remove %d vdisks need:%q\n ", loopCnt, time.Since(t2))

	return nil
}

func case6_Watcher() error{

	agentID := "101"
	hostIp := "10.25.26.46"
	hostname := "agent100"

	common.DeleteAllAgents()
	common.DeleteAllVdisks()

	err := addAgent(agentID, hostIp, hostname)
	if nil != err {
		return err
	}

	watchFunc := etcdIntf.WatchKey()

	_, err = addVdisk("101", "vm_case5", "root/case5/os_vdisk.qcow2")
	if nil != err {
		return err
	}

	go watchFunc("/agents/101/primary_vdisks")

	_, err = addVdisk("101", "vm_case5", "root/case5/os_vdisk2.qcow2")
	if nil != err {
		return err
	}

	return nil
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
	
	err = case6_Watcher()
	if nil != err {
		fmt.Println("case6_Watcher --- Fail, ", err.Error())
	}else{
		fmt.Println("case6_Watcher --- Pass")
	}
}