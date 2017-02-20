package common

import (
	"fmt"
	"sort"
	"encoding/json"
	"strings"
	"errors"
	"vdisk_monitor/src/common/etcdInterface"
)

//Agent structure
//--agent
//    |_agentID
//         |_basicInfo
//         |_primary_vdisks#
//         |_secondary_vdisks
const AGENT_ROOT_NODE string = "agents"
const AGENT_BASICINFO_NODE string = "basicInfo"
const AGENT_PRIMARY_VDISKS_NODE string = "primary_vdisks"
const AGENT_SECONDARY_VDISKS_NODE string = "secondary_vdisks"

//----Begin Agent CRUD----
func GetAgentBasicInfo(agentID string) (AgentBasicInfo, error) {
	
	getValueFunc := etcdIntf.GetKey()
	key := fmt.Sprintf("/%s/%s/%s", AGENT_ROOT_NODE, agentID, AGENT_BASICINFO_NODE)

	value, err := getValueFunc(key)
	if nil != err {
		return AgentBasicInfo{}, err
	}

	var info AgentBasicInfo
	err = json.Unmarshal([]byte(value), &info)
	if nil != err {
		return info, err
	}

	return info, nil
}

func SetAgentBasicInfo(info AgentBasicInfo) error {
	
	setValueFunc := etcdIntf.SetKey()
	key := fmt.Sprintf("/%s/%s/%s", AGENT_ROOT_NODE, info.Id, AGENT_BASICINFO_NODE)

	value, err := json.Marshal(info)
	if nil != err {
		return err
	}

	err = setValueFunc(key, string(value))
	return err
}

func GetPrimaryVdisks(agentID string) ([]string, error){
	
	getValueFunc := etcdIntf.GetKey()
	key := fmt.Sprintf("/%s/%s/%s", AGENT_ROOT_NODE, agentID, AGENT_PRIMARY_VDISKS_NODE)

	value, err := getValueFunc(key)
	if nil != err {
		return []string{}, err
	}

	var vdiskIdList []string

	err = json.Unmarshal([]byte(value), &vdiskIdList)

	return vdiskIdList, err
}

func SetPrimaryVdisks(agentID string, vdiskIdList []string) error {
	
	setValueFunc := etcdIntf.SetKey()
	key := fmt.Sprintf("/%s/%s/%s", AGENT_ROOT_NODE, agentID, AGENT_PRIMARY_VDISKS_NODE)

	value, err := json.Marshal(vdiskIdList)
	if nil != err {
		return err
	}

	err = setValueFunc(key, string(value))
	return err
}

func GetSecondaryVdisks(agentID string) ([]string, error){
	
	getValueFunc := etcdIntf.GetKey()
	key := fmt.Sprintf("/%s/%s/%s", AGENT_ROOT_NODE, agentID, AGENT_SECONDARY_VDISKS_NODE)

	value, err := getValueFunc(key)
	if nil != err {
		return []string{}, err
	}

	var vdiskIdList []string

	err = json.Unmarshal([]byte(value), &vdiskIdList)

	return vdiskIdList, err
}

func SetSecondaryVdisks(agentID string, vdiskIdList []string) error {
	
	setValueFunc := etcdIntf.SetKey()
	key := fmt.Sprintf("/%s/%s/%s", AGENT_ROOT_NODE, agentID, AGENT_SECONDARY_VDISKS_NODE)

	value, err := json.Marshal(vdiskIdList)
	if nil != err {
		return err
	}

	err = setValueFunc(key, string(value))
	return err
}

func SetAgent(agent Agent) error{
	
	err := SetAgentBasicInfo(agent.BasicInfo)
	if nil != err {
		return err
	}

	err = SetPrimaryVdisks(agent.BasicInfo.Id, agent.Primary_vdisks)
	if nil != err {
		return err
	}

	err = SetSecondaryVdisks(agent.BasicInfo.Id, agent.Secondary_vdisks)
	return err
}

func CreateAgent(agent Agent) error {
	
	return SetAgent(agent)
}

func UpdateAgent(agent Agent) error{

	return SetAgent(agent)
}

func DeleteAgent(agentID string) error {

	key := fmt.Sprintf("/%s/%s", AGENT_ROOT_NODE, agentID)

	deleteAgentFunc := etcdIntf.DeleteDirectory()

	err := deleteAgentFunc(key)
	if err != nil {
	 	s := fmt.Sprintf("Delete agent fail! err :%s", err.Error())
	 	return errors.New(s)
	 }

	 return nil
}

func GetAgent(agentID string) (agent Agent, err error) {
	
	agent.BasicInfo, err = GetAgentBasicInfo(agentID)
	if nil != err {
		return agent, err
	}

	agent.Primary_vdisks, err = GetPrimaryVdisks(agentID)
	if nil != err {
		return agent, err
	}

	agent.Secondary_vdisks, err = GetSecondaryVdisks(agentID)
	return agent, nil
}

func GetAgentList() ([]Agent, error){
	
	getAgentListFunc := etcdIntf.GetDirectory()

	agentKeyPaths, err := getAgentListFunc("/agents")
	if nil != err {
		s := fmt.Sprintf("Get agent list fail! Err: %s\n", err.Error())
		return []Agent{}, errors.New(s)
	}

	var agentList []Agent

	for _, path := range agentKeyPaths{
	
		strArr := strings.Split(path, "/")
		agentID := strArr[len(strArr) - 1]

		agent,err := GetAgent(agentID)
		if nil != err {
			continue
		}

		agentList = append(agentList, agent)
	}

	return agentList, err
}

func DeleteAllAgents() error{

	deleteAgentFunc := etcdIntf.DeleteDirectory()

	key := fmt.Sprintf("/%s", AGENT_ROOT_NODE)

	err := deleteAgentFunc(key)
	if err != nil {
	 	s := fmt.Sprintf("Delete all agents fail, err: %s\n", err.Error())
	 	return errors.New(s)
	 }

	 return nil
}


func containsVdisk(vdiskId string, vdiskIdList []string) bool{
	
	for _, id := range vdiskIdList {
	
		if vdiskId == id {
			return true		
		}
	}

	return false
}

func WatchVdiskList(agentID string, bkpType BACKUP_TYPE) (addVdisks []string, rmvVdisks []string, err error){
	
	var key string

	if bkpType == PRIMARY_BACKUP {
		key = fmt.Sprintf("/agents/%s/primary_vdisks")	
	
	}else {
		key = fmt.Sprintf("/agents/%s/secondary_vdisks")			
	}

	watchFunc := etcdIntf.WatchKey()

	currValue, prevValue, err := watchFunc(key)
	if nil != err {
		fmt.Printf("Stop watching %s, Err:%s", agentID, err.Error())
		return
	}

	var currVdiskList, prevVdiskList []string
	json.Unmarshal([]byte(currValue), &currVdiskList)
	json.Unmarshal([]byte(prevValue), &prevVdiskList)

	sort.Strings(currVdiskList)
	sort.Strings(prevVdiskList)

	var found bool

	//check new vdiskId
	for _, vdiskId := range currVdiskList {

		found = containsVdisk(vdiskId, prevVdiskList)
		if false == found {
			addVdisks = append(addVdisks, vdiskId)
		}
	}

	//check removed vdiskId
	for _, vdiskId := range prevVdiskList {

		found = containsVdisk(vdiskId, currVdiskList)
		if false == found {
			rmvVdisks = append(rmvVdisks, vdiskId)
		}
	}

	return
}
