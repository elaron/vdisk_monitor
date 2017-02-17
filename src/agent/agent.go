package main

import (
	"fmt"
	"net"
	"time"
	"errors"
	"flag"
	"vdisk_monitor/src/common/dbConn"
)

type AgentConfig struct {
	AgentId		string
	PortRange [2]uint32
	CurrPort uint32
	TcpServerPort uint32
	HostIp	string
}

var g_agentConfig AgentConfig

func main() {

	initAgentConfig()
	sendRegisteAgentMsg()
	
	//go heartbeatToMds()
	go sendAddVdiskMsgToMds()
	go setuplistener()

	for {
        time.Sleep(60 * time.Second)
	}
	
}

func initAgentConfig() {

    port := flag.Int("port", 3333, "Tcp server port")
    ip := flag.String("ip", "127.0.0.1", "server ip")
    agentId := flag.String("id", "agent100", "agent Id")
 
    flag.Parse()

	g_agentConfig.AgentId = *agentId
	g_agentConfig.TcpServerPort = uint32(*port)
	g_agentConfig.HostIp = *ip

	g_agentConfig.PortRange[0] = 20000
	g_agentConfig.PortRange[1] = 40000

	g_agentConfig.CurrPort = g_agentConfig.PortRange[0]
}

func getAvailablePort(num uint8) ([]uint32, error) {

	portRange := g_agentConfig.PortRange[1] - g_agentConfig.PortRange[0]
		
	var i uint32
	var ports []uint32
	for i = 0; i < portRange; i++ {
			
		tryPort := g_agentConfig.PortRange[0] + ((g_agentConfig.CurrPort + i)%portRange)
			
		addr := fmt.Sprintf("localhost:%d", tryPort)
		netListen, err := net.Listen("tcp", addr)
		if nil != err {
			continue	
		}

		defer netListen.Close()	

		ports = append(ports, tryPort)
		num--

		if 0 == num {
			g_agentConfig.CurrPort = tryPort + 1
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

	syncInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, common.PRIMARY_BACKUP)
	if nil != err {
		return
	}

	syncInfo.Tcp_server_port = ports
	syncInfo.LastHeartBeatTime = time.Now().String()

	common.SetVdiskBackupDaemonInfo(vdiskId, syncInfo, common.PRIMARY_BACKUP)
}
