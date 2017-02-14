package main
	
import (
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
	"errors"
	"strconv"
)

func formatInt32(n int32) string {
    return strconv.FormatInt(int64(n), 10)
}

//Agent structure
//--agent
//    |_agentID
//         |_basicInfo
//         |_primary_vdisks
//         |_secondary_vdisks

//----Begin Agent CRUD----
const (
	AGENT_BASIC_INFO = iota
	AGENT_PRIMARY_VDISKS
	AGENT_SECONDARY_VDISKS
	AGGENT_SUB_NODE_TYPE_BUTT
)
type AGENT_SUB_NODE_TYPE int32

func getAgentRootKey(agentID int32) string{

	key := bytes.Buffer{}
	key.WriteString("/agents/")
	key.WriteString(formatInt32(agentID))
	key.WriteString("/")
	
	return key.String()
}

func getAgentSubNodeNames() []string {
	return []string{"basicInfo", "primary_vdisks", "secondary_vdisks"}
}

func getAgentSubNodesPaths(agentID int32) []string{
	
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

func  getAgentSubNodeValues(agent Agent) ([AGGENT_SUB_NODE_TYPE_BUTT]string, error){
	
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

func setAgent(f func() (func(key string, value string) error), agent Agent) error{
	
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

func createAgent(agent Agent) error {
	
	return setAgent(createKey, agent)
}

func updateAgent(agent Agent) error{

	return setAgent(updateKey, agent)
}

func deleteAgent(agentID int32) error {

	agentRootKey 	:= getAgentRootKey(agentID)	
	deleteAgentFunc := deleteDirectory()

	err := deleteAgentFunc(agentRootKey)
	if err != nil {
	 	fmt.Println("Delete agent fail")
	 }

	 return err
}

func getAgent(agentID int32) (Agent, error) {
	
	getAgentNodeValueFunc := getKey()
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

func getAgentList() ([]Agent, error){
	
	getAgentListFunc := getDirectory()

	agentKeyPaths, err := getAgentListFunc("/agents")
	if nil != err {
		fmt.Println("Get agent list fail! Err:", err.Error())
		return []Agent{}, err
	}

	var agentList []Agent

	for _, path := range agentKeyPaths{
	
		subStrs := strings.Split(path, "/")

		agentID, err := strconv.ParseInt(subStrs[len(subStrs) - 1], 10, 32)
		if err != nil {
    		continue
		}

		agent,err := getAgent(int32(agentID))
		if nil != err {
			continue
		}

		agentList = append(agentList, agent)
	}

	return agentList, err
}

func deleteAllAgents() error{
	deleteAgentFunc := deleteDirectory()

	err := deleteAgentFunc("/agents")
	if err != nil {
	 	fmt.Println("Delete all agents fail")
	 }

	 return err
}
//vdisk CRUD