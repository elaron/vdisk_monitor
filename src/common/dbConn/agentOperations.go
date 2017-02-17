package common

import (
	"fmt"
	"encoding/json"
	"bytes"
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

//----Begin Agent CRUD----
const (
	AGENT_BASIC_INFO = iota
	AGENT_PRIMARY_VDISKS
	AGENT_SECONDARY_VDISKS
	AGGENT_SUB_NODE_TYPE_BUTT
)
type AGENT_SUB_NODE_TYPE int32

func getAgentRootKey(agentID string) string{

	key := bytes.Buffer{}
	key.WriteString("/agents/")
	key.WriteString(agentID)
	key.WriteString("/")
	
	return key.String()
}

func getAgentSubNodeNames() []string {
	return []string{"basicInfo", "primary_vdisks", "secondary_vdisks"}
}

func getAgentSubNodesPaths(agentID string) []string{
	
	agentRootKey := getAgentRootKey(agentID)
	subNodeNames := getAgentSubNodeNames()

	var path bytes.Buffer
	var subNodePaths []string
	
	for _, name := range subNodeNames {
		path.Reset()
		path.WriteString(agentRootKey)
		path.WriteString(name)

		subNodePaths = append(subNodePaths, path.String())
	}

	return subNodePaths
}

func getAgentSubNodeValues(agent Agent) ([AGGENT_SUB_NODE_TYPE_BUTT]string, error){
	
	var b 				[]byte
	var errs 			[AGGENT_SUB_NODE_TYPE_BUTT]error
	var subNodeValues 	[AGGENT_SUB_NODE_TYPE_BUTT]string
	
	b, errs[AGENT_BASIC_INFO] = json.Marshal(agent.BasicInfo)
	subNodeValues[AGENT_BASIC_INFO] = string(b)

	b, errs[AGENT_PRIMARY_VDISKS] = json.Marshal(agent.Primary_vdisks)
	subNodeValues[AGENT_PRIMARY_VDISKS] = string(b)

	b, errs[AGENT_SECONDARY_VDISKS] = json.Marshal(agent.Secondary_vdisks)
	subNodeValues[AGENT_SECONDARY_VDISKS] = string(b)

	for _, e := range errs {
		if e != nil {
			fmt.Printf("Encode %s fail", e.Error())
			return [AGGENT_SUB_NODE_TYPE_BUTT]string{}, e
		}
	}

	return subNodeValues, nil
}

func SetAgent(f func() (func(key string, value string) error), agent Agent) error{
	
	nodeValues, err := getAgentSubNodeValues(agent)
	if err != nil {
		fmt.Println("Update agent fail, err:", err.Error())
		return err
	}

	nodeNames 		:= getAgentSubNodeNames()	
	subNodePaths 	:= getAgentSubNodesPaths(agent.BasicInfo.Id)	
	addAgentFunc 	:= f()

	for i, name := range nodeNames {
	
		errMsg := addAgentFunc(subNodePaths[i], nodeValues[i])
		if errMsg != nil {
			fmt.Printf("Set agent's %s fail. Key:%s Value:%s", name, subNodePaths[i], nodeValues[i])
			return errMsg
		}
	}

	return nil
}

func CreateAgent(agent Agent) error {
	
	return SetAgent(etcdIntf.CreateKey, agent)
}

func UpdateAgent(agent Agent) error{

	return SetAgent(etcdIntf.UpdateKey, agent)
}

func DeleteAgent(agentID string) error {

	agentRootKey 	:= getAgentRootKey(agentID)	
	deleteAgentFunc := etcdIntf.DeleteDirectory()

	err := deleteAgentFunc(agentRootKey)
	if err != nil {
	 	s := fmt.Sprintf("Delete agent fail! err :%s", err.Error())
	 	return errors.New(s)
	 }

	 return nil
}

func GetAgent(agentID string) (Agent, error) {
	
	getAgentNodeValueFunc := etcdIntf.GetKey()
	subNodePaths := getAgentSubNodesPaths(agentID)

	var value [AGGENT_SUB_NODE_TYPE_BUTT]string

	for i, path := range subNodePaths {
		value[i], _ = getAgentNodeValueFunc(path)
		//fmt.Printf("Value[%d]:%s\n", i, value[i])	//just for debug
	}

	var agent Agent
	var errs [AGGENT_SUB_NODE_TYPE_BUTT]error 

	if 0 != len(value[AGENT_BASIC_INFO]) {
		errs[AGENT_BASIC_INFO] = json.Unmarshal([]byte(value[AGENT_BASIC_INFO]), &agent.BasicInfo)
	
	}else{
		return Agent{}, errors.New("Key is non-exist")
	}

	if 0 != len(value[AGENT_PRIMARY_VDISKS]) {
		errs[AGENT_PRIMARY_VDISKS] = json.Unmarshal([]byte(value[AGENT_PRIMARY_VDISKS]), &agent.Primary_vdisks)	
	}

	if 0 != len(value[AGENT_SECONDARY_VDISKS]) {
		errs[AGENT_SECONDARY_VDISKS] = json.Unmarshal([]byte(value[AGENT_SECONDARY_VDISKS]), &agent.Secondary_vdisks)
	}

	for i, e := range errs {
		if e != nil {
			fmt.Printf("Unmarshal %s's value fail! Err: %s. \n", subNodePaths[i], e.Error())
			return Agent{}, e
		}
	}

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

	err := deleteAgentFunc("/agents")
	if err != nil {
	 	s := fmt.Sprintf("Delete all agents fail, err: %s\n", err.Error())
	 	return errors.New(s)
	 }

	 return nil
}
