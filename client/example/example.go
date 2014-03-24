package main

import (
    "crypto/rand"
    "fmt"
    "time"

    "github.com/abhi-bit/goq/client"
)

func main() {
    var nodeList client.Node
    //Connect to cluster manager
    nodeList, _  = client.Connect("http://localhost:9091/nodes")
    fmt.Printf("%#v\n", nodeList)

    for {
        client.EnqueueData(randString(6))
        time.Sleep(time.Second)
    }

}

func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}
