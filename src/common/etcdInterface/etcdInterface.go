package etcdIntf

import (
	"log"
	"time"
	"fmt"
	"errors"
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

func CreateKey() func(key string, value string) error{
	
	keyApi := getKeysAPI()

	return func (key string, value string) error {
		
		// Create is an alias for Set w/ PrevExist=false
		//Create(ctx context.Context, key, value string) (*Response, error)
		
		//resp, err := keyApi.Create(context.Background(), key, value)
		_, err := keyApi.Create(context.Background(), key, value)
		if err != nil {
			log.Println(err)

		}/* else {
			log.Printf("Create is done. Metadata is %q\n", resp)
		}*/

		return err
	}
}

func DeleteKey() func(key string) error{

	// Delete removes a Node identified by the given key, optionally destroying
	// all of its children as well. The caller may define a set of required
	// conditions in an DeleteOptions object.
	//Delete(ctx context.Context, key string, opts *DeleteOptions) (*Response, error)

	keyApi := getKeysAPI()

	return func (key string) error {
		
		//resp, err := keyApi.Delete(context.Background(), key, nil)
		_, err := keyApi.Delete(context.Background(), key, nil)

		if err != nil {
			s := fmt.Sprintf("Delete key fail, err: %s", err.Error())
			return errors.New(s)		
		}/*else{
			log.Printf("Delete is done, Metadata is %q\n", resp)
		}*/

		return nil
	}
}

func UpdateKey() func(key string, value string) error{
	
	keyApi := getKeysAPI()

	return func (key string, value string) error {

		// Update is an alias for Set w/ PrevExist=true
		//Update(ctx context.Context, key, value string) (*Response, error)
		
		//resp, err := keyApi.Update(context.Background(), key, value)
		_, err := keyApi.Update(context.Background(), key, value)
		
		if err != nil {
			log.Println(err)
			fmt.Println("Update key fail", err.Error())
		
		} /*else {
			// print common key info
			log.Printf("Update is done. Metadata is %q\n", resp)
		}*/

		return err
	}
}

func SetKey() func(key string, value string) error{
	
	keyApi := getKeysAPI()

	return func (key string, value string) error {

	// Set assigns a new value to a Node identified by a given key. The caller
	// may define a set of conditions in the SetOptions. If SetOptions.Dir=true
	// then value is ignored.
	//Set(ctx context.Context, key, value string, opts *SetOptions) (*Response, error)

		_, err := keyApi.Set(context.Background(), key, value, nil)
		
		if err != nil {
			log.Println(err)
			fmt.Println("Set key fail", err.Error())
		
		} /*else {
			// print common key info
			log.Printf("Update is done. Metadata is %q\n", resp)
		}*/

		return err
	}
}

func GetKey() func(key string) (string, error){
	
	keyApi := getKeysAPI()

	return func (key string) (string, error) {

		// Get retrieves a set of Nodes from etcd
		//Get(ctx context.Context, key string, opts *GetOptions) (*Response, error)

		resp, err := keyApi.Get(context.Background(), key, nil)
		
		if err != nil {
			//log.Fatal(err)
			s := fmt.Sprintf("Etcd - Get key fail, err : %s\n", err.Error())
			return "", errors.New(s)

		}/*else {
			log.Printf("Get is done. Metadata is %q\n", resp)
			log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
		}*/

		return string(resp.Node.Value), err
	}
}

func DeleteDirectory() func (key string) error{
	
	keyApi := getKeysAPI()
	
	opt := &client.DeleteOptions{
			PrevIndex: 0,
			Recursive: true,
			Dir: true,
		}

	return func (key string) error {
		
		_, err := keyApi.Delete(context.Background(), key, opt)

		if err != nil {
			//log.Fatal(err)
			fmt.Println("Delete directory fail, ", err.Error())
		
		}

		return err
	}
}

func GetDirectory() (func (key string) ([]string, error)) {
	
	keyApi := getKeysAPI()

	opt := &client.GetOptions{
		Recursive : true,
		Sort : true,
		Quorum : true,
	}

	return func (key string) ([]string, error) {

		var values []string

		resp, err := keyApi.Get(context.Background(), key, opt)
		if nil != err {
			
			s := fmt.Sprintf("Get directory fail, err: %s\n", err.Error())
			return []string{}, errors.New(s)

		}else{
			
			if resp.Node.Dir != true {
				values = append(values, string(resp.Node.Value))

			}else{
								
				for _,node := range resp.Node.Nodes {
				
					//fmt.Printf("%d key:%s value:%s\n", i, node.Key, node.Value)
					values = append(values, string(node.Key))
				}
			}
		}

		return values, nil
	}
}

func WatchKey()(func (key string) (currentValue string, prevValue string ,err error)) {
		
	keyApi := getKeysAPI()

	return func (key string) (currentValue string, prevValue string ,err error) {
		
		resp, err := keyApi.Get(context.Background(), key, nil)
		if nil != err {
			s := fmt.Sprintf("Get key(%s) fail, cannot watch non-exist key!", key)
			err = errors.New(s)
			return
		}

		//Watcher(key string, opts *WatcherOptions) Watcher
		watcher := keyApi.Watcher(key, nil)

		resp, err = watcher.Next(context.Background())
		if err != nil {
			return 
		}
/*
		fmt.Println("Watch response:", resp)
		fmt.Printf("action:%s \n value:%s \n prevValue:%s\n",
			resp.Action,
			resp.Node.Value,
			resp.PrevNode.Value)
*/
		currentValue = string(resp.Node.Value)
		prevValue = string(resp.PrevNode.Value)

		return 
	}
}