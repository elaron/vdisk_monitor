package main

import (
	"fmt"
	"net"
	"time"
	"errors"
)

type AgentConfig struct {
	portRange [2]uint32
	currPort uint32
}

var g_agentConfig AgentConfig

func main() {

	initAgentConfig()
	sendRegisteAgentMsg()
	
	//go heartbeatToMds()
	go sendAddVdiskMsgToMds()

	for {
        time.Sleep(60 * time.Second)
	}
	
}

func initAgentConfig() {

	g_agentConfig.portRange[0] = 20000
	g_agentConfig.portRange[1] = 40000

	g_agentConfig.currPort = g_agentConfig.portRange[0]
}

func getAvailablePort(num uint8) ([]uint32, error) {

	portRange := g_agentConfig.portRange[1] - g_agentConfig.portRange[0]
		
	var i uint32
	var ports []uint32
	for i = 0; i < portRange; i++ {
			
		tryPort := g_agentConfig.portRange[0] + ((g_agentConfig.currPort + i)%portRange)
			
		addr := fmt.Sprintf("localhost:%d", tryPort)
		netListen, err := net.Listen("tcp", addr)
		if nil != err {
			continue	
		}

		defer netListen.Close()	

		ports = append(ports, tryPort)
		num--

		if 0 == num {
			g_agentConfig.currPort = tryPort + 1
			return ports, nil
		}
	}

	return []uint32{}, errors.New("Get ports fail")	
}

func startOriginator(vdiskId string) {
	
	ports, err := getAvailablePort(4)
	if nil != err {
		fmt.Println(err.Error())
	}

	fmt.Println(vdiskId, ports)
}
