package main

func main() {

	sendRegisteAgentMsg()
	
	//go heartbeatToMds()
	go sendAddVdiskMsgToMds()

	for {

	}
	
}
