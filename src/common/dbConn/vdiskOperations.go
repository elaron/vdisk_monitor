package common
	
import (
	"fmt"
	"encoding/json"
	"errors"
	"strconv"
	"vdisk_monitor/src/common/etcdInterface"
)

//Vdisk structure
//--vdisks
//    |_vdiskID
//    |    |_vmInfo
//    |    |_backups
//    |    		|_primary_bkp
//    |    				|_backupInfo
//    |    				|_daemonInfo
//    |    		|_secondary_bkp
//    |    				|_backupInfo
//    |    				|_daemonInfo
//    |_waitToRemove

const VDISK_ROOT_NODE string = "vdisks"
const VDISK_VM_INFO_NODE string = "vmInfo"
const VDISK_BACKUPS_NODE string = "backups"
const VDISK_PRIMARY_BKP_NODE string = "primary_bkp"
const VDISK_SECONDARY_BKP_NODE string = "secondary_bkp"
const VDISK_BACKUP_INFO_NODE string = "backupInfo"
const VDISK_DAEMON_INFO_NODE string = "daemonInfo"

const VDISK_WAIT_TO_REMOVE string = "waitToRemove"

func formatInt32(n int32) string {
    return strconv.FormatInt(int64(n), 10)
}

//vdisk CRUD
func getVdiskSubNodeKey(vdiskId string, subNode string) string{
	
	key := fmt.Sprintf("/%s/%s/%s", VDISK_ROOT_NODE, vdiskId, subNode)
	return key
}

func getVdiskVmInfoKey(vdiskId string) string{
	str := getVdiskSubNodeKey(vdiskId, VDISK_VM_INFO_NODE)
	return str
}

func getVdiskBackupKey(vdiskId string, backupType BACKUP_TYPE) string{

	var subNode string
	switch backupType {
	
	case PRIMARY_BACKUP:
		subNode = fmt.Sprintf("/%s/%s/%s", VDISK_BACKUPS_NODE, VDISK_PRIMARY_BKP_NODE, VDISK_BACKUP_INFO_NODE)

	case SECONDARY_BACKUP:
		subNode = fmt.Sprintf("/%s/%s/%s", VDISK_BACKUPS_NODE, VDISK_SECONDARY_BKP_NODE, VDISK_BACKUP_INFO_NODE)
	
	default:
		fmt.Printf("Invalid backupType:%d\n", backupType)
		subNode = ""
	}
	
	key := getVdiskSubNodeKey(vdiskId, subNode)
	return key
}

func getVdiskBackupDaemonInfoKey(vdiskId string, backupType BACKUP_TYPE) string{
	
	var subNode string

	switch backupType {

	case PRIMARY_BACKUP:
		subNode = fmt.Sprintf("/%s/%s/%s", VDISK_BACKUPS_NODE, VDISK_PRIMARY_BKP_NODE, VDISK_DAEMON_INFO_NODE)

	case SECONDARY_BACKUP:
		subNode = fmt.Sprintf("/%s/%s/%s", VDISK_BACKUPS_NODE, VDISK_SECONDARY_BKP_NODE, VDISK_DAEMON_INFO_NODE)

	default:
		subNode = ""
	}

	key := getVdiskSubNodeKey(vdiskId, subNode)
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

	fmt.Printf("Delete vdisk(%s) success!\n", vdiskId)

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

func GetWaitToRemoveVdiskList() ([]string, error){

	key := fmt.Sprintf("/%s/%s", VDISK_ROOT_NODE, VDISK_WAIT_TO_REMOVE)	
	
	getListFunc := etcdIntf.GetKey()

	value, err := getListFunc(key)
	if nil != err {
		return []string{}, err
	}

	var list []string
	err = json.Unmarshal([]byte(value), &list)

	return list,err
}

func AddVdiskToRemoveList(vdiskId string) error {
	
	list, _ := GetWaitToRemoveVdiskList()
	list = append(list, vdiskId)

	
	key := fmt.Sprintf("/%s/%s", VDISK_ROOT_NODE, VDISK_WAIT_TO_REMOVE)	
	value, err := json.Marshal(list)
	if nil != err {
		return err
	}

	setListFunc := etcdIntf.SetKey()
	err = setListFunc(key, string(value))
	
	return err
}

func RemoveVdiskIdInList(vdiskId string, vdiskIdList []string) []string{

	if len(vdiskIdList) == 0 {
		return []string{}
	}

	var rmvIdx int	
	for index, id := range vdiskIdList {
		if vdiskId == id {
			rmvIdx = index
			break
		}
	}

	newList := append(vdiskIdList[:rmvIdx], vdiskIdList[rmvIdx+1:]...)
	return newList
}

func EraseVdiskFromRemoveList(vdiskId string) error {
	
	list, err := GetWaitToRemoveVdiskList()
	if nil != err {
		return err
	}

	list = RemoveVdiskIdInList(vdiskId, list)
	
	key := fmt.Sprintf("/%s/%s", VDISK_ROOT_NODE, VDISK_WAIT_TO_REMOVE)	
	value, err := json.Marshal(list)
	if nil != err {
		return err
	}

	setListFunc := etcdIntf.SetKey()
	err = setListFunc(key, string(value))
	
	return err
}
