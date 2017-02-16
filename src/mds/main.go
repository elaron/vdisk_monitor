package main

import (
    "fmt"
    "time"
)

func main() {
	
	go setuplistener()

	for {
		fmt.Println("I am in main now!")
        time.Sleep(60 * time.Second)
	}
}
