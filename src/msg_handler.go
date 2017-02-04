package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	fmt.Println("hello")

	agent := Agent{
		HostIp:     "10.25.26.46",
		Hostname:   "agent100",
		Id:         100,
		State:      ACTIVE,
		Orig_vdisk: []string{"5d33e3b0-9e39-4ecc-9847-f3fc7a97b0c5", "47f17840-f800-4881-a8f0-584354ab24b8"},
		Term_vdisk: []string{"47f17840-f800-4881-a8f0-584354ab24b8", "0e0c92f1-1692-4d48-8af5-bbceae73caa4"},
	}

	//fmt.Println(agent)

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
