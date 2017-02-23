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

func startTerminator(vdiskId string){

	ports, err := getAvailablePort(4)
	if nil != err {
		fmt.Println(err.Error())
	}

	fmt.Println(common.SECONDARY_BACKUP, vdiskId, ports)

	origInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, common.PRIMARY_BACKUP)
	if nil != err {
		return
	}

	termInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, common.SECONDARY_BACKUP)
	if nil != err {
		return
	}

	fmt.Println("Originator Tcp_server_port:", origInfo.Tcp_server_port)

	//Use Orig TCP port and Term tcp port to start terminator

	termInfo.Tcp_server_port = ports
	termInfo.LastHeartBeatTime = time.Now().String()
	termInfo.State = common.ACTIVE

	common.SetVdiskBackupDaemonInfo(vdiskId, termInfo, common.SECONDARY_BACKUP)
}

func watchTerminatorState(vdiskId string, timer *time.Timer) {
	
	state, err := common.WatchSyncDaemonState(vdiskId, common.SECONDARY_BACKUP)
	if err != nil {
		fmt.Printf("Watch terminator(%s) fail!\n", vdiskId)
		return
	}

	if common.ACTIVE == state {
		fmt.Printf("Add vdisk(%s) success!\n", vdiskId)			
		
		//Start terminator success ,so there's no need to remove this vdisk
		timer.Stop()
	}
}

func startOriginator(vdiskId string) {

	//get available Tcp server port
	ports, err := getAvailablePort(4)
	if nil != err {
		fmt.Println(err.Error())
	}

	fmt.Println(common.PRIMARY_BACKUP, vdiskId, ports)


	//start sync originator
	//TODO

	syncInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, common.PRIMARY_BACKUP)
	if nil != err {
		return
	}

	syncInfo.Tcp_server_port = ports
	syncInfo.LastHeartBeatTime = time.Now().String()
	syncInfo.State = common.ACTIVE
	
	timer := time.NewTimer(5 * time.Second)
	go func() {
        <- timer.C

        //If time's up and terminator is still not active, 
        //then we think this ADD_VDISK operation is fail
        fmt.Println("Time out. Add vdisk fail!")
        g_buff.rmvVdisks <- vdiskId
    }()

    //watch sync terminator
	go watchTerminatorState(vdiskId, timer)

	//update sync orignator daemonInfo on etcd
	err = common.SetVdiskBackupDaemonInfo(vdiskId, syncInfo, common.PRIMARY_BACKUP)
	if err != nil {
		fmt.Printf("Set originator info fail, Err :%s \n", err.Error())
	}
}

func startSync(vdiskId string, bkpType common.BACKUP_TYPE) {
	
	if common.PRIMARY_BACKUP == bkpType {
		go startOriginator(vdiskId)

	} else if common.SECONDARY_BACKUP == bkpType{
		go startTerminator(vdiskId)
	
	}else{
		fmt.Printf("Unrecognize bkpType:%d, start sync for %s fail!\n", 
			bkpType,
			vdiskId)
	}

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
			
			time.Sleep(30 * time.Second)
			continue
		}

		for _,vdiskId := range addVdisks {
			fmt.Printf("Start new vdisk %d %s %s\n", bkpType, vdiskId, time.Now())
			startSync(vdiskId, bkpType)
		}

		for _,vdiskId := range rmvVdisks {
			fmt.Printf("Remove vdisk %d %s %s\n", bkpType, vdiskId, time.Now())
			removeSync(vdiskId, bkpType)
		}
	}
}


