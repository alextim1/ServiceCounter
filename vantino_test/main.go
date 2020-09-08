package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"sync"
	"log"
)

// Request structure represents an object receiving from the Logger.
type Request struct{
	TimeStamp  	string
	Ip 			string
	URL			string
}

// IpCounter structure wraps the map of unique ip addresses.
type IpCounter struct{
	uniqueIp sync.Map
}

func (ct IpCounter) post(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
    var t Request
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
	}
	
	ct.uniqueIp.Store(t.Ip, "visited")
  }

func (ct IpCounter) get(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s unique ip addresses: %d\n", ct.uniqueIp.Range)
}

func main() {
	counter := IpCounter{}
	http.HandleFunc("/logs", counter.post)
	http.HandleFunc("/metrics", counter.get)
	
	// goroutine to lunc a server for Prometheus requests.
	go func(){
		log.Fatal(http.ListenAndServe("localhost:9102"))
	}
	
	log.Fatal(http.ListenAndServe("localhost:5000", nil))
}
