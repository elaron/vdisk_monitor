package main

import (
	"log"
	"time"
	"fmt"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func getKeysAPI () client.KeysAPI{
	
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

	return kapi
}

func setKey() func(key string, value string) error{
	
	keyApi := getKeysAPI()

	return func (key string, value string) error {

		resp, err := keyApi.Set(context.Background(), key, value, nil)
		
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
	
	keyApi := getKeysAPI()

	return func (key string) (string, error) {

		resp, err := keyApi.Get(context.Background(), key, nil)

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

func deleteKey() func(key string) error{
	//Delete(ctx context.Context, key string, opts *DeleteOptions) (*Response, error)

	keyApi := getKeysAPI()

	return func (key string) error {
		
		resp, err := keyApi.Delete(context.Background(), key, nil)

		if err != nil {
			log.Fatal(err)
			fmt.Println("Delete key fail, ", err.Error())
		
		}else{
			log.Printf("Delete is done, Metadata is %q\n", resp)
		}

		return err
	}
}