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
	ResidentAgentID 	string
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
	VmState 			VM_STATE
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

	if len(agent.Primary_vdisks) == 0 {
		return false
	}

	for _, vdiskId := range agent.Primary_vdisks {

		bkpInfo,_ := getVdiskBackupInfo(vdiskId, PRIMARY_BACKUP)
		if vdiskPath == bkpInfo.Path {
			return true
		}
	}

	return false
}

func addVdisk(agentID string, vmId string, path string) error{
	
	agent,err := getAgent(agentID)
	if nil != err {
		s := fmt.Sprintf("Agent(%s) is not exist!")
		return errors.New(s)
	}

	vdiskExist := isVdiskExistOnAgent(path, agent)
	if true == vdiskExist {
		s := fmt.Sprintf("Vdisk(%s) already exist on agent!\n", path)
		return errors.New(s)
	}

	var vdisk Vdisk

	vdisk.Id = genUUID()
	vdisk.VmInfo.VmId = vmId
	vdisk.VmInfo.VmState = RUNNING

	vdisk.Backups[PRIMARY_BACKUP].BackupInfo.ResidentAgentID = agentID
	vdisk.Backups[PRIMARY_BACKUP].BackupInfo.Path = path

	vdisk.Backups[PRIMARY_BACKUP].SyncDaemonInfo.SyncType = ORIGINATOR

	err = createVdisk(vdisk)
	if nil != err {
		s := fmt.Sprintf("Create vdisk fail! Add vdisk(%s) on agent(%s) fail!", path, agentID)
		return errors.New(s)
	}
	agent.Primary_vdisks = append(agent.Primary_vdisks, vdisk.Id)

	err = updateAgent(agent)
	if nil != err {
		s := fmt.Sprintf("Update agent(%s) fail, add vdisk(%s) fail!", agentID, path)
		return errors.New(s)
	}
	
	return err
}

func getVdiskId(agentID string, path string) (string, error){
	
	agent,err := getAgent(agentID)
	if nil != err {
		s := fmt.Sprintf("Get vdiskId(path=%s) fail, because get agent(%s) fail", path, agentID)
		return "", errors.New(s)
	}

	for _, vdiskId := range agent.Primary_vdisks {

		bkpInfo,_ := getVdiskBackupInfo(vdiskId, PRIMARY_BACKUP)
		
		if path == bkpInfo.Path {
			return vdiskId, nil
		}
	}

	s := fmt.Sprintf("Can't find vdisk(%s) on agent(%s)", path, agentID)
	return "", errors.New(s)
}

func removeVdisk(vdiskID string, agentID string, path string) error{

	var rmvVdiskID string
	var err error

	if len(vdiskID) == 0 {
		
		rmvVdiskID, err = getVdiskId(agentID, path)
		if nil != err{
			return err
		}
	}else{
		rmvVdiskID = vdiskID
	}

	err = deleteVdisk(rmvVdiskID)
	return err
}