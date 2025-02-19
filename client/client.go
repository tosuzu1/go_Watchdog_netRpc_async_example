package main

import (
	"go_netRpc_async_test/data"
	"log"
	"net/rpc"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func isError_Log(err error) bool {
	if err != nil {
		log.Printf("Error : %v\n", err)
		return true
	} else {
		return false
	}
}

func isError_Fatal(err error) bool {
	// fatal crash if error
	if err != nil {
		log.Fatalf("Error : %v\n", err)
		return true
	} else {
		return false
	}
}

func main() {
	var configStruct data.Config_struct
	ClientHeartBeatChan := make(chan *rpc.Call, 100)

	// get name of config file from command line argument or default
	var configFilename string
	switch len(os.Args) {
	case 1:
		configFilename = "../data/local_config.yaml" // default confg file name
	case 2:
		configFilename = os.Args[1]
	default:
		log.Fatal("Too many arguments provided")
	}

	// This section handle the yaml file
	// Get file name of config
	configFile, err := os.ReadFile(configFilename)
	isError_Fatal(err)

	// Unmasrhall yaml
	err = yaml.Unmarshal(configFile, &configStruct)
	isError_Fatal(err)

	log.Printf("Client info: %v\n", configStruct.Client)
	log.Printf("Server info: %v\n", configStruct.Server)

	for {
		serverRpc, err := rpc.DialHTTP("tcp", configStruct.Server.Ipv4_address+":"+configStruct.Server.Port)
		if isError_Log(err) {
			// we eounctered a error, continue to search for another server
			time.Sleep(2 * time.Second)
			continue
		}
		// RPC called async
		var resp time.Duration
		theCall := serverRpc.Go("ServerAPI.ClientHeartBeat", 0, &resp, ClientHeartBeatChan)

		select {
		case chanData := <-ClientHeartBeatChan:
			if !isError_Log(chanData.Error) {
				log.Printf("Heartbeat sent to server : WD elaspe = %v\n", chanData.Reply)
				log.Printf("DEBUG chandData = %v\n", chanData)
			}
		case <-time.After(30 * time.Millisecond):
			log.Printf("Timeout!\n")
		}

		log.Printf("DEBUG: theCall = %v\n", *theCall)

		time.Sleep(2 * time.Second)
	}

}
