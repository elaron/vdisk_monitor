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
	PeerAgentId string
	PortRange [2]uint32
	CurrPort uint32
	TcpServerPort uint32
	HostIp	string
}

type BufferChannel struct {
	rmvVdisks chan string
	addVdisks chan string
}

var g_agentConfig AgentConfig
var g_buff BufferChannel

func main() {

	initAgentConfig()
	initAgentBuffer()

	sendRegisteAgentMsg()

	for {
        time.Sleep(60 * time.Second)
	}
	
}

func initAgentConfig() {

    port := flag.Int("port", 3333, "Tcp server port")
    ip := flag.String("ip", "127.0.0.1", "server ip")
    agentId := flag.String("id", "agent100", "agent Id")
    peerAgentId := flag.String("peerId", "agent200", "peer-agent Id")
 
    flag.Parse()

	g_agentConfig.AgentId = *agentId
	g_agentConfig.PeerAgentId = *peerAgentId

	g_agentConfig.TcpServerPort = uint32(*port)
	g_agentConfig.HostIp = *ip

	g_agentConfig.PortRange[0] = 20000
	g_agentConfig.PortRange[1] = 40000

	g_agentConfig.CurrPort = g_agentConfig.PortRange[0]
}

func initAgentBuffer() {
	g_buff.rmvVdisks = make(chan string, 10)
	g_buff.addVdisks = make(chan string, 10)
}

func runAgent() {
	go watchVdiskChange(g_agentConfig.AgentId, common.PRIMARY_BACKUP)
	go watchVdiskChange(g_agentConfig.AgentId, common.SECONDARY_BACKUP)

	go mock_sendAddVdiskMsgToMds()
	go sendRemoveVdiskMsgToMds()
	go heartbeatToMds()
	go setuplistener()
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

func startSync(vdiskId string, bkpType common.BACKUP_TYPE) {
	
	ports, err := getAvailablePort(4)
	if nil != err {
		fmt.Println(err.Error())
	}

	fmt.Println(bkpType, vdiskId, ports)

	syncInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, bkpType)
	if nil != err {
		return
	}

	syncInfo.Tcp_server_port = ports
	syncInfo.LastHeartBeatTime = time.Now().String()

	common.SetVdiskBackupDaemonInfo(vdiskId, syncInfo, bkpType)
}

func removeSync(vdiskId string, bkpType common.BACKUP_TYPE) {
	
	syncInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, bkpType)
	if nil != err {
		return
	}

	//send commond to stop sync daemon

	syncInfo.State = common.REMOVED
	syncInfo.Tcp_server_port = []uint32{}

	common.SetVdiskBackupDaemonInfo(vdiskId, syncInfo, bkpType)
}

func watchVdiskChange(agentId string, bkpType common.BACKUP_TYPE) {
	
	for {
		addVdisks, rmvVdisks, err := common.WatchVdiskList(agentId, bkpType)

		if err != nil {
			list := []string{"primary_vdisks", "secondary_vdisks"}
			fmt.Printf("Stop watching %s vdisk list\n", list[bkpType])
			continue
		}

		for _,vdiskId := range addVdisks {
			fmt.Printf("Start new vdisk %s\n", vdiskId)
			startSync(vdiskId, bkpType)
		}

		for _,vdiskId := range rmvVdisks {
			fmt.Printf("Remove vdisk %s\n", vdiskId)
			removeSync(vdiskId, bkpType)
		}
	}
}


