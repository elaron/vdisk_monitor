package main

func main() {
	
	go heartbeatToMds()
	go sendAddVdiskMsgToMds()

	for {

	}
	
}
