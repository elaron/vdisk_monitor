package main

import "time"

func main() {

	sendRegisteAgentMsg()
	
	//go heartbeatToMds()
	go sendAddVdiskMsgToMds()

	for {
        time.Sleep(60 * time.Second)
	}
	
}
