package main

import (
	"time"
	"errors"
)

type AgentConfig struct {
	portRange [2]uint32
}

var g_agentConfig AgentConfig

func main() {

	sendRegisteAgentMsg()
	
	//go heartbeatToMds()
	go sendAddVdiskMsgToMds()

	for {
        time.Sleep(60 * time.Second)
	}
	
}

func getAvailablePort() (ports [4]uint32, err error){
	return [4]uint32{}, errors.New("")
}

func startOriginator() {
	
}
