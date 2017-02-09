package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	//"strings"
)

const (
	ACTIVE = iota
	LOSS_CONN
	STATE_TYPE_BUTT
)

type AGENT_STATE int32
type Agent struct {
	HostIp     string
	Hostname   string
	Id         int32
	State      AGENT_STATE
	Orig_vdisk []string
	Term_vdisk []string
}

func genUUID() string {

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println("genUUID fail!")
	}
	return string(uuid[0:36])
}

func main() {
	const MAX_DISK_NUM = 128
	var orig_disks []string
	var term_disks []string

	for i := 0; i < MAX_DISK_NUM; i++ {
		var uuid1 string = genUUID()
		orig_disks = append(orig_disks, uuid1)
		uuid2 := genUUID()
		term_disks = append(term_disks, uuid2)
	}

	//fmt.Println(strings.Join(term_disks, ", "))

	agent := Agent{
		HostIp:     "10.25.26.46",
		Hostname:   "agent100",
		Id:         100,
		State:      ACTIVE,
		Orig_vdisk: orig_disks,
		Term_vdisk: term_disks,
	}

	b, err := json.Marshal(agent)
	if nil != err {
		fmt.Println("encode to json fail!")
	} else {
		fmt.Println("json-body:", string(b))
	}
}
func msg_dispatch(msg string) bool {
	return true
}
