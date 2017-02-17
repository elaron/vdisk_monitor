package main

import (
	"fmt"
	"os/exec"
	"errors"
	"vdisk_monitor/src/common/dbConn"
)

func genUUID() string {

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println("genUUID fail!")
	}
	return string(uuid[0:36])
}

func isAgentExist(agentList []common.Agent, agent common.Agent) (bool, string){

	for _, item := range agentList {

		if agent.BasicInfo.HostIp == item.BasicInfo.HostIp {
			return true, "Duplicate agent IP"

		}else if agent.BasicInfo.Id == item.BasicInfo.Id {
			return true, "Duplicate agent ID"

		}
	}
	return false, ""
}

func addAgent(agentId string, ip string, hostname string) error{

	agent := common.Agent{
		BasicInfo: common.AgentBasicInfo {
			HostIp:     ip,
			Hostname:   hostname,
			Id:         agentId,
			State: common.ACTIVE,
		},
	}

	agentList, _ := common.GetAgentList()

	agentExist, msg := isAgentExist(agentList, agent)
	if true == agentExist {
		s := fmt.Sprintf("Add agent fail, because of %s", msg)
		return errors.New(s)
	}

	err := common.CreateAgent(agent)
	if nil != err {
		s := fmt.Sprintf("Create agent fail, Err: %s", err.Error())
		return errors.New(s)
	}

	return nil
}

func removeAgent(agentID string) error{
	
	//funcName, file, line, _ := runtime.Caller(0)
	
	_, err := common.GetAgent(agentID)
	if nil != err {
		return errors.New("Agnet dose'n exist")
	}

	err = common.DeleteAgent(agentID)
	if nil != err {
		s := fmt.Sprintf("Delete agent fail, Err:%s", err.Error())
		return errors.New(s)
	}

	return nil
}

func isVdiskExistOnAgent(vdiskPath string, agent common.Agent) bool{

	if len(agent.Primary_vdisks) == 0 {
		return false
	}

	for _, vdiskId := range agent.Primary_vdisks {

		bkpInfo,_ := common.GetVdiskBackupInfo(vdiskId, common.PRIMARY_BACKUP)
		if vdiskPath == bkpInfo.Path {
			return true
		}
	}

	return false
}

func addVdisk(agentID string, vmId string, path string) (string, error){
	
	agent, err := common.GetAgent(agentID)
	if nil != err {
		s := fmt.Sprintf("Agent(%s) is not exist!")
		return "", errors.New(s)
	}

	vdiskId, _ := getVdiskId(agentID, path)
	if len(vdiskId) != 0 {
		s := fmt.Sprintf("Vdisk(%s) already exist on agent!\n", path)
		return vdiskId, errors.New(s)
	}

	var vdisk common.Vdisk

	vdisk.Id = genUUID()
	vdisk.VmInfo.VmId = vmId
	vdisk.VmInfo.VmState = common.RUNNING

	vdisk.Backups[common.PRIMARY_BACKUP].BackupInfo.ResidentAgentID = agentID
	vdisk.Backups[common.PRIMARY_BACKUP].BackupInfo.Path = path

	vdisk.Backups[common.PRIMARY_BACKUP].SyncDaemonInfo.SyncType = common.ORIGINATOR

	err = common.CreateVdisk(vdisk)
	if nil != err {
		s := fmt.Sprintf("Create vdisk fail! Add vdisk(%s) on agent(%s) fail!", path, agentID)
		return vdisk.Id, errors.New(s)
	}
	agent.Primary_vdisks = append(agent.Primary_vdisks, vdisk.Id)

	err = common.UpdateAgent(agent)
	if nil != err {
		s := fmt.Sprintf("Update agent(%s) fail, add vdisk(%s) fail!", agentID, path)
		return vdisk.Id, errors.New(s)
	}
	
	return vdisk.Id, nil
}

func getVdiskId(agentID string, path string) (string, error){
	
	agent,err := common.GetAgent(agentID)
	if nil != err {
		s := fmt.Sprintf("Get vdiskId(path=%s) fail, because get agent(%s) fail", path, agentID)
		return "", errors.New(s)
	}

	for _, vdiskId := range agent.Primary_vdisks {

		bkpInfo,_ := common.GetVdiskBackupInfo(vdiskId, common.PRIMARY_BACKUP)
		
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

	err = common.DeleteVdisk(rmvVdiskID)
	return err
}