package main

import (
	"go_netRpc_async_test/data"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func checkErrorLog(err error) bool {
	if err != nil {
		log.Printf("Error : %v\n", err)
		return true
	} else {
		return false
	}
}

func checkErrorFatal(err error) bool {
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
	checkErrorFatal(err)

	// Unmasrhall yaml
	err = yaml.Unmarshal(configFile, &configStruct)
	checkErrorFatal(err)

	log.Printf("Client info: %v\n", configStruct.Client)
	log.Printf("Server info: %v\n", configStruct.Server)

	// start watchdog api
	thisapi := new(ServerAPI)
	log.Println("Registering host API")
	rpc.Register(thisapi)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", configStruct.Server.Ipv4_address+":"+configStruct.Server.Port)
	checkErrorFatal(err)
	go http.Serve(listener, nil)

	http.HandleFunc("/api", thisapi.HandleJSONHeartBeat)
	go http.ListenAndServe(":8080", nil)

	thisapi.Init(5000) // 5 second should be good enought

	// start a go routine for the timer, this will count down the watchdog timer
	go thisapi.StartWatchDog()

	for {
		// loop, check to make sure watchdog is alive, upon timeout exit program
		if thisapi.IsWatchTimeout() {
			break
		}
		log.Printf("WD elasped time %v\n", thisapi.getElapsedTime())
		time.Sleep(10 * time.Millisecond)
	}

	log.Printf("Watchdog timedout!\n")
}
