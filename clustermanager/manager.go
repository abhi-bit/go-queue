package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/couchbaselabs/clog"
)

var (
	address string
	port    int
	logPath string
	hosts   string
	nodes   = make(map[string]NodeStatus)
)

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&address, "address", "", "Address to listen on, Default is to all")
	flag.IntVar(&port, "port", 9091, "Port to listen on. Default is 8091")
	flag.StringVar(&logPath, "path", "manager", "cluster manager logging dir")
	flag.StringVar(&hosts, "host", "localhost:11311", "nodes to manage")
	flag.Parse()

}

type NodeStatus struct {
	status  bool
	retries int
}

func Nodes(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	onlineNodes := make([]string, 0)
	for node, _ := range nodes {
		onlineNodes = append(onlineNodes, node)
	}

	oNodes, _ := json.Marshal(onlineNodes)

	fmt.Fprintf(w, fmt.Sprintf("{\"nodes\":%s}", oNodes))
}

func main() {

	log.Printf("listening on %s:%d\n", address, port)
	log.Printf("cluster manager Path: %s\n", logPath)

	for _, host := range strings.Split(hosts, ",") {
		serverURL := "http://" + host

		resp, err := http.Get(serverURL)

		if err != nil {
			clog.Error(err)
			nodes[serverURL] = NodeStatus{status: false, retries: 0}
			break
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if string(body) == "1" {
			nodes[serverURL] = NodeStatus{status: true, retries: 0}
		} else {
			nodes[serverURL] = NodeStatus{status: false, retries: 0}
		}
	}

	fmt.Printf("%#v\n", nodes)

	//Polling nodes, needs cleanup
	go func() {
		for {
			for node, _ := range nodes {
				resp, err := http.Get(node)
				if err != nil {
					clog.Error(err)
					retryCount := nodes[node].retries + 1

					if retryCount <= 3 {
						nodes[node] = NodeStatus{status: false, retries: retryCount}
					} else {
						delete(nodes, node)
					}

					break
				}

				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if string(body) == "1" {
					nodes[node] = NodeStatus{status: true, retries: 0}
				} else {
					retryCount := nodes[node].retries + 1
					if retryCount <= 3 {
						nodes[node] = NodeStatus{status: false, retries: retryCount}
					} else {
						delete(nodes, node)
					}
				}
			}
			fmt.Printf("%#v\n", nodes)
			time.Sleep(time.Second)
		}
	}()

	http.HandleFunc("/nodes", Nodes)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
	if err != nil {
		log.Fatalf("Failed to start cluster manager: %v", err)
	}
}
