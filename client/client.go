package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Node struct {
	Nodes []string `json:"nodes"`
}

type Data struct {
	success bool
	data    []string `json:"data"`
	message string
}

var (
	address     string
	port        int
	nodes       Node
	workCounter int
	workMux     sync.Mutex
)

func Connect(clusterURL string) (Node, error) {

	resp, err := http.Get(clusterURL)
	if err != nil {
		return Node{}, err
	}

	//Need to handle this in-case cluster manager dies
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &nodes)
	return nodes, nil

}

func EnqueueData(nodes Node, data string) error {
	values := make(url.Values)
	values.Set("data", data)

	shard := findShard()
	DBServerURL := nodes.Nodes[shard] + "/enqueue"

	req, err := http.PostForm(DBServerURL, values)
	if err != nil {
		return err
	}

	resp, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	defer req.Body.Close()
	log.Printf("data: %s, shard: %s, resp: %s\n", data, DBServerURL, resp)

	return nil
}

func DequeueData(DBServerURL string) (Data, error) {

	ServerURL := DBServerURL + "/dequeue"

	resp, err := http.Get(ServerURL)
	if err != nil {
		return Data{}, err
	}

	defer resp.Body.Close()

	var dData Data

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &dData)

	return dData, nil
}

func findShard() int {
	workMux.Lock()
	workCounter++
	workMux.Unlock()

	nodeCount := len(nodes.Nodes)
	return (workCounter % nodeCount)
}
