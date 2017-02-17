package common
	
import (
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
	"errors"
	"strconv"
	"vdisk_monitor/src/common/etcdInterface"
)

func formatInt32(n int32) string {
    return strconv.FormatInt(int64(n), 10)
}
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

//vdisk CRUD
func getVdiskSubNodeKey(vdiskId string, subNode string) string{
	
	key := bytes.Buffer{}
	
	key.WriteString("/vdisks/")
	key.WriteString(vdiskId)
	key.WriteString(subNode)
	
	return key.String()
}

func getVdiskVmInfoKey(vdiskId string) string{
	str := getVdiskSubNodeKey(vdiskId, "/vmInfo")
	return str
}

func getVdiskBackupKey(vdiskId string, backupType BACKUP_TYPE) string{

	var key string
	switch backupType {
	
	case PRIMARY_BACKUP:
		key = getVdiskSubNodeKey(vdiskId, "/backups/primary_bkp/backupInfo")

	case SECONDARY_BACKUP:
		key = getVdiskSubNodeKey(vdiskId, "/backups/secondary_bkp/backupInfo")
	
	default:
		fmt.Printf("Invalid backupType:%d\n", backupType)
		key = ""
	}

	return key
}

func getVdiskBackupDaemonInfoKey(vdiskId string, backupType BACKUP_TYPE) string{
	
	var key string

	switch backupType {

	case PRIMARY_BACKUP:
		key = getVdiskSubNodeKey(vdiskId, "/backups/primary_bkp/daemonInfo")

	case SECONDARY_BACKUP:
		key = getVdiskSubNodeKey(vdiskId, "/backups/secondary_bkp/daemonInfo")

	default:
		key = ""
	}

	return key
}

func setNodeValue(key string, obj interface{}) error{
	
	value, err := json.Marshal(obj)
	if nil != err {
		s := fmt.Sprintf("Set key:%s value fail.\n", key)
		return errors.New(s)
	}

	setValueFunc := etcdIntf.SetKey()

	err = setValueFunc(key, string(value))
	if nil != err {
		return err
	}

	return nil	
}

func SetVdiskVmInfo(vdiskId string, info VmInfomation) error{
	
	key := getVdiskVmInfoKey(vdiskId)
	err := setNodeValue(key, info)

	return err
}

func SetVdiskBackupInfo(vdiskId string, bkp VdiskBackupInfo, bkpType BACKUP_TYPE) error{
	
	key := getVdiskBackupKey(vdiskId, bkpType)
	err := setNodeValue(key, bkp)

	return err
}

func SetVdiskBackupDaemonInfo(vdiskId string, daemonInfo SyncDaemon, bkpType BACKUP_TYPE) error{

	key := getVdiskBackupDaemonInfoKey(vdiskId, bkpType)
	err := setNodeValue(key, daemonInfo)

	return err
}

func getNodeOjb(key string) (obj interface{}, err error) {

	getValueFunc := etcdIntf.GetKey()

	value,err := getValueFunc(key)
	if nil != err {
		return 
	}

	err = json.Unmarshal([]byte(value), &obj)

	return
}

func GetVdiskVmInfo(vdiskId string) (info VmInfomation, err error) {

	key := getVdiskVmInfoKey(vdiskId)
	getValueFunc := etcdIntf.GetKey()

	value, err := getValueFunc(key)
	if nil != err {
		return 
	}

	err = json.Unmarshal([]byte(value), &info)
	return
}

func GetVdiskBackupInfo(vdiskId string, bkpType BACKUP_TYPE) (info VdiskBackupInfo, err error){

	key := getVdiskBackupKey(vdiskId, bkpType)
	getValueFunc := etcdIntf.GetKey()

	value, err := getValueFunc(key)
	if nil != err {
		return 
	}

	err = json.Unmarshal([]byte(value), &info)
	return
}

func GetVdiskBackupDaemonInfo(vdiskId string, bkpType BACKUP_TYPE) (info SyncDaemon, err error){
	
	key := getVdiskBackupDaemonInfoKey(vdiskId, bkpType)
	getValueFunc := etcdIntf.GetKey()

	value, err := getValueFunc(key)
	if nil != err {
		return 
	}

	err = json.Unmarshal([]byte(value), &info)
	return
}

func CreateVdisk(vdisk Vdisk) error{
	
	var err error

	err = SetVdiskVmInfo(vdisk.Id, vdisk.VmInfo)
	if nil != err {
		return err
	}

	var bkpType BACKUP_TYPE = PRIMARY_BACKUP

	for ; bkpType < BACKUP_TYPE_BUTT; bkpType++ {
	
		err = SetVdiskBackupInfo(vdisk.Id, vdisk.Backups[bkpType].BackupInfo, bkpType)
		if nil != err {
			return err
		}
		
		err = SetVdiskBackupDaemonInfo(vdisk.Id, vdisk.Backups[bkpType].SyncDaemonInfo, bkpType)
		if nil != err {
			return err
		}		
	}

	return nil
}

func DeleteVdisk(vdiskId string) error{
	
	deleteKeyFunc := etcdIntf.DeleteDirectory()
	key := fmt.Sprintf("/vdisks/%s", vdiskId)

	err := deleteKeyFunc(key)
	if nil != err {
		s := fmt.Sprintf("Delete vdisk fail! err :%s", err.Error())
		return errors.New(s)
	}

	return nil
}

func DeleteAllVdisks() error{

	deleteAgentFunc := etcdIntf.DeleteDirectory()

	err := deleteAgentFunc("/vdisks")
	if err != nil {
	 	s := fmt.Sprintf("Delete all vdisks fail, err: %s\n", err.Error())
	 	return errors.New(s)
	 }

	 return nil
}
