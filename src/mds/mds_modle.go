package main

import (
	"fmt"
	"os/exec"
	//"strings"
)

const (
	ACTIVE = iota
	LOSS_CONN
	STATE_TYPE_BUTT
)

type AGENT_STATE int32
type vdiskList []string

type Agent struct {
	HostIp     string
	Hostname   string
	Id         int32
	State      AGENT_STATE
}

func genUUID() string {

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println("genUUID fail!")
	}
	return string(uuid[0:36])
}

func isAgentExist(agentList []Agent, agent Agent) bool{

	for _, item := range agentList {

		if agent.HostIp == item.HostIp {
			fmt.Println("HostIp exist", item.HostIp)
			return true

		}else if agent.Id == item.Id {
			fmt.Println("agentID exist:", item.Id)
			return true

		}
	}
	return false
}

