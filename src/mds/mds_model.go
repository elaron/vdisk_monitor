package main

import (
	"fmt"
	"os/exec"
	//"strings"
)

//define enums
const (
	ACTIVE = iota
	LOSS_CONN
	DAEMON_STATE_TYPE_BUTT
)
type DAEMON_STATE_TYPE int32

const (
	RUNNING = iota
	PAUSE
	STOP
	MIGRATE
	VM_STATE_TYPE_BUTT
)
type VM_STATE int32

const (
	ORIGINATOR = iota
	TERMINATOR
	SYNC_TYPE_BUTT
)
type SYNC_TYPE int32

const (
	NORMAL_SYNC = iota
	INCREASE_SYNC
	FULL_SYNC
	BACKUP_STATE_BUTT
)
type BACKUP_STATE int32

//define basic structure
type AgentID 		int32
type vdiskList 		[]string
type Agent struct {
	HostIp string
	Hostname string
	Id AgentID
	State DAEMON_STATE_TYPE
}

type SyncDaemon struct {
	SyncType 			SYNC_TYPE
	Tcp_server_port 	int32
	LastWriteSeq		int64
	State 				DAEMON_STATE_TYPE
	LastHeartBeatTime	int64
}

type VdiskBackup struct {
	ResidentAgentID 	AgentID
	Path 				string
	Size 				int64
	BackupStatus 		BACKUP_STATE
	SyncPercent 		int32
	SyncDaemonInfo		SyncDaemon
}

type Vdisk struct {
	Id 					string
	VmId 				string
	Vmstate 			VM_STATE
	Backups 			[]VdiskBackup
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

//func addAgent(agent Agent) error{}

func isVdiskExistOnAgent(path string){

}

