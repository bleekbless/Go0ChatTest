package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var portNumber *string
var configFile *string
var debug bool
var config Config

func botsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	err := readRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func statHandler(w http.ResponseWriter, r *http.Request) {
	if debug {
		fmt.Printf("[%v] %v\n", r.Method, r.URL)
	}

	printStat(w)

	if debug {
		printStat(os.Stdout)
	}
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	if debug {
		fmt.Printf("[%v] %v\n", r.Method, r.URL)
	}

	resetStat(config)
	printStat(w)

	if debug {
		printStat(os.Stdout)
	}
}

func main() {
	portNumber = flag.String("port", "8877", "Port number to use for connection")
	configFile = flag.String("conf", "config.json", "Path to config file")
	debugFlag := flag.Bool("debug", false, "Write debug information")
	flag.Parse()
	debug = *debugFlag

	readConfig(&config)
	initStat(config)

	fmt.Println("I'm ready!")
	address := "http://127.0.0.1:" + *portNumber
	fmt.Println(address + "/stat - use it for watching statistic during a test")
	fmt.Println(address + "/reset - use it to reset statistic between tests")

	http.HandleFunc("/stat", statHandler)
	http.HandleFunc("/reset", resetHandler)
	http.HandleFunc("/", botsHandler)

	log.Fatal(http.ListenAndServe(":"+*portNumber, nil))
}
