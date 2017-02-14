package main

import (
	"fmt"
	"runtime"
	"os/exec"
	"errors"
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
type AgentBasicInfo struct {
	HostIp 		string
	Hostname 	string
	Id 			string
	State 		DAEMON_STATE_TYPE
}

type Agent struct {
	BasicInfo 		AgentBasicInfo
	Primary_vdisks 	[]string
	Secondary_vdisks []string
}

type SyncDaemon struct {
	SyncType			SYNC_TYPE
	Tcp_server_port		int32
	LastWriteSeq		int64
	State 				DAEMON_STATE_TYPE
	LastHeartBeatTime	int64
}

type VdiskBackupInfo struct {
	ResidentAgentID 	int32
	Path 				string
	Size 				int64
	BackupStatus 		BACKUP_STATE
	SyncPercent 		int32
}

type VdiskBackup struct {
	BackupInfo 		VdiskBackupInfo
	SyncDaemonInfo	SyncDaemon
}

const (
	PRIMARY_BACKUP = iota
	SECONDARY_BACKUP
	BACKUP_TYPE_BUTT
)
type BACKUP_TYPE int32

type VmInfomation struct {
	VmId 				string
	Vmstate 			VM_STATE
}

type Vdisk struct {
	Id 			string
	VmInfo 		VmInfomation
	Backups 	[BACKUP_TYPE_BUTT]VdiskBackup
 }

func genUUID() string {

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println("genUUID fail!")
	}
	return string(uuid[0:36])
}

func isAgentExist(agentList []Agent, agent Agent) (bool, string){

	for _, item := range agentList {

		if agent.BasicInfo.HostIp == item.BasicInfo.HostIp {
			return true, "Duplicate agent IP"

		}else if agent.BasicInfo.Id == item.BasicInfo.Id {
			return true, "Duplicate agent ID"

		}
	}
	return false, ""
}

func addAgent(agent Agent) (bool, string){

	agentList, err := getAgentList()

	agentExist, errMsg := isAgentExist(agentList, agent)
	if true == agentExist {
		return false, errMsg
	}

	err = createAgent(agent)
	if nil != err {
		return false, "Add agent fail"
	}

	return true, ""
}

func removeAgent(agentID string) (bool, string){
	
	funcName, file, line, _ := runtime.Caller(0)
	
	_, err := getAgent(agentID)
	if nil != err {
			return false, "Agnet not exist"
	}

	err = deleteAgent(agentID)
	if nil != err {
		fmt.Printf("[%s][%s][%d]Remove agent fail!, Err%q\n ", runtime.FuncForPC(funcName).Name(), file, line, err.Error())
		return false, "Delete agent fail"
	}

	return true,""
}

func isVdiskExistOnAgent(vdiskPath string, agent Agent) bool{
/*
	for _, vdiskId := range agent.Primary_vdisks {


	}*/

	return false
}

func addVdisk(agentID string, vmId string, path string) error{
	
	agent,err := getAgent(agentID)
	if nil != err {
		fmt.Printf("Can't add vdisk on nonexistent agent(id=%d)\n", agentID)
		return errors.New("Agent not exist!")
	}

	vdiskExist := isVdiskExistOnAgent(path, agent)
	if true == vdiskExist {
		errMsg := fmt.Sprintf("Vdisk(%s) already exist on agent!\n", path)
		return errors.New(errMsg)
	}

	return nil
}