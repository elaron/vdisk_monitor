package main

import (
	"log"
	"time"
	"fmt"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func setKey() func(key string, value string) error{
	
	cfg := client.Config{
	
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	
	if err != nil {
		log.Fatal(err)
		fmt.Println("new client faile", err.Error())
	}
	
	kapi := client.NewKeysAPI(c)

	return func (key string, value string) error {

		resp, err := kapi.Set(context.Background(), key, value, nil)
		
		if err != nil {
			log.Fatal(err)
			fmt.Println("set key fail", err.Error())
		
		} else {
			// print common key info
			log.Printf("Set is done. Metadata is %q\n", resp)
		}

		return err
	}
}


func getKey() func(key string) (string, error){
	
	cfg := client.Config{
	
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	
	if err != nil {
		log.Fatal(err)
		fmt.Println("new client faile", err.Error())
	}
	
	kapi := client.NewKeysAPI(c)

	return func (key string) (string, error) {

		resp, err := kapi.Get(context.Background(), key, nil)

		if err != nil {
			log.Fatal(err)
			fmt.Println("get key fail", err.Error())

		} else {
			log.Printf("Get is done. Metadata is %q\n", resp)
			log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
		}

		return string(resp.Node.Value), err
	}
}