# go-queue

go-queue is distributed persistent queue implemented in Golang.


## Features

- Supported API (`enequeue`, `dequeue`, `alldocs`, `statistics`, `version`)
- All communication between client, DB and cluster manager over HTTP

## Installation

```
go get github.com/abhi-bit/go-queue
go build *.go
```

## Queue Configuration

One can run multipe instances of qdb on same box, just need to make sure their ports are unique. To start the queue, different options:

- **port** - Port qdb should listen on, defaults to 11311

- **sync** - Synchronise to LevelDB on every write, defaults to true

- **path** - Path to the LevelDB datavase directory.

- **address** - Address that this instance of qdb should listen on.


## Cluster Manager Configuration:

- **address** - Address that cluster manager will be running on

- **host** - Queue instances to manage, defaults to localhost:11311

- **port** - Port cluster manager should listen on, defaults to 9091


## Client Library:

- API supported `Enqueue` and `Dequeue`

- Sample code using client library available under client/example.


