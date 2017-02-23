package common

//define enums
const (
	UNAVAILABLE = iota
	ACTIVE 
	LOSS_CONN
	REMOVED
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

const (
	PRIMARY_BACKUP = iota
	SECONDARY_BACKUP
	BACKUP_TYPE_BUTT
)
type BACKUP_TYPE int32

//define basic structure
type AgentBasicInfo struct {
	HostIp 		string
	Hostname 	string
	Id 			string
	PeerAgentId	string
	State 		DAEMON_STATE_TYPE
	TcpServerPort uint32
}

type Agent struct {
	BasicInfo 		AgentBasicInfo
	Primary_vdisks 	[]string
	Secondary_vdisks []string
}

type SyncDaemon struct {
	SyncType			SYNC_TYPE
	Tcp_server_port		[]uint32
	LastWriteSeq		int64
	State 				DAEMON_STATE_TYPE
	LastHeartBeatTime	string
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


type VmInfomation struct {
	VmId 				string
	VmState 			VM_STATE
}

type Vdisk struct {
	Id 			string
	VmInfo 		VmInfomation
	Backups 	[BACKUP_TYPE_BUTT]VdiskBackup
 }