package main

import (
	"fmt"
)

func case1_findExistedAgent() (bool, string){

	//create agentList
	var agentList []Agent

	agent1  := Agent{
		HostIp:     "10.25.26.46",
		Hostname:   "aaa",
		Id:         100,
		State:      ACTIVE,
	}

	agent2  := Agent{
		HostIp:     "10.25.26.46",
		Hostname:   "bbb",
		Id:         101,
		State:      ACTIVE,
	}

	agentList = append(agentList, agent1)
	agentList = append(agentList, agent2)

	var agentTest Agent
	var rslt bool

	//test duplicate agentID
	agentTest = Agent{
		HostIp:     "192.25.26.46",
		Hostname:   "lalalala",
		Id:         100,
		State:      ACTIVE,
	}

	rslt = isAgentExist(agentList, agentTest)
	if rslt != true {
		fmt.Println("case1_findExistedAgent-testFail")
		return false, "duplicate agentID"
	}

	//test duplicate IP
	agentTest = Agent{
		HostIp:     "10.25.26.46",
		Hostname:   "lulululu",
		Id:         200,
		State:      ACTIVE,
	}

	rslt = isAgentExist(agentList, agentTest)
	if rslt != true {
		fmt.Println("case1_findExistedAgent-testFail")
		return false, "duplicate hostIP"
	}

	return true, ""
}

func main() {
	
	rslt, failReason := case1_findExistedAgent()
	if rslt != true {
		fmt.Println("case1_findExistedAgent Fail! failReason:", failReason)
	}else{
		fmt.Println("case1_findExistedAgent PASS")
	}

	etcdClient()

}