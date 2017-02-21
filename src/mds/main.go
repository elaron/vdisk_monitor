package main

import (
    "fmt"
    "time"
)

func main() {
	
	go setuplistener()
    setupTimers()

	for {
		fmt.Println("I am in main now!")
        time.Sleep(60 * time.Second)
	}
}

func setupTimers() {
    
    ticker := time.NewTicker(time.Second * 1)
    go func() {
        for range ticker.C {
            checkWaitToRemoveVdisk()
        }
    }()
}
