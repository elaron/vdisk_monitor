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

	if false == agentExist {

		err := common.CreateAgent(agent)
		if nil != err {
			s := fmt.Sprintf("Create agent fail, Err: %s", err.Error())
			return errors.New(s)
		}

	}else{

		info, err := common.GetAgentBasicInfo(agentId)
		
		ok := (nil == err) && (info.HostIp == ip) && (info.Id == agentId)
		if false == ok {
			s := fmt.Sprintf("Add agent fail, because of %s", msg)
			return errors.New(s)				
		}
	}

	return nil
}

func setPeerAgent(agentID string, peerAgentId string) error {
	
	info, err := common.GetAgentBasicInfo(agentID)
	if nil != err {
		return err
	}

	info.PeerAgentId = peerAgentId

	err = common.SetAgentBasicInfo(info)
	return err
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
	
	//check agent/peerAgent existence and vdisk's duplication
	agent, err := common.GetAgent(agentID)
	if nil != err {
		s := fmt.Sprintf("Agent(%s) is not exist!", agentID)
		return "", errors.New(s)
	}

	peerAgent, err := common.GetAgent(agent.BasicInfo.PeerAgentId)
	if nil != err {
		s := fmt.Sprintf("Peer Agent(%s) is not exist, add vdisk on Agent(%s) fail!", 
			agent.BasicInfo.PeerAgentId, 
			agentID)
		return "", errors.New(s)
	}

	vdiskId, _ := getVdiskId(agentID, path)
	if len(vdiskId) != 0 {
		s := fmt.Sprintf("Vdisk(%s) already exist on agent!\n", path)
		return vdiskId, errors.New(s)
	}

	//create new vdisk info node
	var vdisk common.Vdisk

	vdisk.Id = genUUID()
	vdisk.VmInfo.VmId = vmId
	vdisk.VmInfo.VmState = common.RUNNING

	vdisk.Backups[common.PRIMARY_BACKUP].BackupInfo.ResidentAgentID = agentID
	vdisk.Backups[common.PRIMARY_BACKUP].BackupInfo.Path = path
	vdisk.Backups[common.PRIMARY_BACKUP].SyncDaemonInfo.SyncType = common.ORIGINATOR

	vdisk.Backups[common.SECONDARY_BACKUP].BackupInfo.ResidentAgentID = agent.BasicInfo.PeerAgentId
	vdisk.Backups[common.SECONDARY_BACKUP].BackupInfo.Path = path
	vdisk.Backups[common.SECONDARY_BACKUP].SyncDaemonInfo.SyncType = common.TERMINATOR

	err = common.CreateVdisk(vdisk)
	if nil != err {
		s := fmt.Sprintf("Create vdisk fail! Add vdisk(%s) on agent(%s) fail!", path, agentID)
		return vdisk.Id, errors.New(s)
	}

	go func () {

		vdiskId := vdisk.Id
	
		//fmt.Printf("Start watching originator daemonInfo %s\n", time.Now())
		state, err := common.WatchSyncDaemonState(vdiskId, common.PRIMARY_BACKUP)
		if err != nil {
			fmt.Printf("Watch originator(%s) fail!\n", vdiskId)
			return
		}

		if common.ACTIVE != state {
			fmt.Printf("Originator of %s is not runing!\n", vdiskId)		
			return
		}

		//append new vdiskId to peerAgent's secondary vdiskList
		newSecondaryVdiskList := append(peerAgent.Secondary_vdisks, vdisk.Id)
	
		err = common.SetSecondaryVdisks(agent.BasicInfo.PeerAgentId, newSecondaryVdiskList)
		if nil != err {
			fmt.Printf("Update peerAgent(%s) fail, add vdisk(%s) fail!", agentID, path)
			return
		}
	}()

	//append new vdiskId into agent's primary vdiskList and peerAgent's secondary vdiskList
	newPrimaryVdiskList := append(agent.Primary_vdisks, vdisk.Id)
	
	err = common.SetPrimaryVdisks(agentID, newPrimaryVdiskList)
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

func removeVdisk(vdiskId string) error {

	//remove vdiskID from PrimaryAgent's vdiskList
	primaryBkp, err := common.GetVdiskBackupInfo(vdiskId, common.PRIMARY_BACKUP)
	if err != nil {
		return err
	}

	primaryVdisks, err := common.GetPrimaryVdisks(primaryBkp.ResidentAgentID)
	if nil != err {
		return err
	}
	newList := common.RemoveVdiskIdInList(vdiskId, primaryVdisks)

	err = common.SetPrimaryVdisks(primaryBkp.ResidentAgentID, newList)
	if nil != err {
		return err
	}

	//remove vdiskID from PrimaryAgent's vdiskList
	secondaryBkp, err := common.GetVdiskBackupInfo(vdiskId, common.SECONDARY_BACKUP)
	if err != nil {
		return err
	}

	secondaryVdisks, err := common.GetSecondaryVdisks(secondaryBkp.ResidentAgentID)
	if nil != err {
		return err
	}
	newList = common.RemoveVdiskIdInList(vdiskId, secondaryVdisks)

	err = common.SetSecondaryVdisks(secondaryBkp.ResidentAgentID, newList)
	if nil != err {
		return err
	}

	err = common.AddVdiskToRemoveList(vdiskId)
	if nil != err {
		s := fmt.Sprintf("Add vdiskId to WaitToRemoveVdiskList fail, Err:%s ", err.Error())
		return errors.New(s)
	}

	return nil
}

func removeVdisk2(vdiskID string, agentID string, path string) error{

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

func checkWaitToRemoveVdisk() {
	
	list, err := common.GetWaitToRemoveVdiskList()

	if nil != err {
		return
	}

	for _, vdiskId := range list {
		
		origInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, common.PRIMARY_BACKUP)
		if nil != err {
			continue
		}

		termInfo, err := common.GetVdiskBackupDaemonInfo(vdiskId, common.SECONDARY_BACKUP)
		if nil != err {
			continue
		}

		canBeRemoved := (origInfo.State == common.REMOVED || origInfo.State == common.UNAVAILABLE) && (termInfo.State == common.REMOVED || termInfo.State == common.UNAVAILABLE)
		if true == canBeRemoved {
			common.DeleteVdisk(vdiskId)
			common.EraseVdiskFromRemoveList(vdiskId)
		}
	}
}