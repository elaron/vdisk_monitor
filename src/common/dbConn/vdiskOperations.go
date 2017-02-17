package common
	
import (
	"fmt"
	"encoding/json"
	"bytes"
	"errors"
	"strconv"
	"vdisk_monitor/src/common/etcdInterface"
)

func formatInt32(n int32) string {
    return strconv.FormatInt(int64(n), 10)
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
