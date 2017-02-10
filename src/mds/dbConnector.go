package main
	
import (
	"fmt"
	"encoding/json"
	"bytes"
	"strconv"
)

func formatInt32(n int32) string {
    return strconv.FormatInt(int64(n), 10)
}

//Agent CRUD
func addAgent(agent Agent) error{

	b, err := json.Marshal(agent.BasicInfo)

	if nil != err {
		fmt.Println("encode to json fail!")
	} 

	key := bytes.Buffer{}

	key.WriteString("/agents/")
	key.WriteString(formatInt32(agent.BasicInfo.Id))
	key.WriteString("/BasicInfo")

	addNewAgent := createKey()
	err = addNewAgent(key.String(), string(b))
	if err != nil {
		fmt.Println("Add agent fail")
	}

	return err
}

func deleteAgent(agentID int32) error {
	
	deleteAgentOp := deleteDirectory()
	
	key := bytes.Buffer{}

	key.WriteString("/agents/")
	key.WriteString(formatInt32(agentID))

	err := deleteAgentOp(key.String())
	if err != nil {
	 	fmt.Println("Delete agent fail")
	 }

	 return err
}

func getAgent(agentID int32) (Agent, error) {
	
	getAgentInfo := getKey()

	key := bytes.Buffer{}

	key.WriteString("/agents/")
	key.WriteString(formatInt32(agentID))
	key.WriteString("/BasicInfo")

	value, err := getAgentInfo(key.String())
	if err != nil {
		fmt.Println("Get agent fail", err.Error())
		return Agent{}, err
	}

	var agent Agent
	err = json.Unmarshal([]byte(value), &agent.BasicInfo)

	fmt.Println(agent)

	return agent, err
}

//vdisk CRUD