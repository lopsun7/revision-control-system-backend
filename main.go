package main

import (
	"github.com/hashicorp/raft"
	"log"
	"strconv"
	"time"
)

func main() {
	nodeCount := 3
	raftPortStart := 5000
	nodes := make([]*raft.Raft, nodeCount)
	transports := make([]raft.Transport, nodeCount)
	ids := make([]raft.ServerID, nodeCount)

	for i := 0; i < nodeCount; i++ {
		nodeID := "node" + strconv.Itoa(i+1)
		raftPort := raftPortStart + i
		nodes[i], transports[i], ids[i] = setupRaft(nodeID, raftPort)
		time.Sleep(1 * time.Second) // Ensure nodes have time to start properly
	}

	// Bootstrap the cluster with the first node
	config := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      ids[0],
				Address: transports[0].LocalAddr(),
			},
		},
	}
	future := nodes[0].BootstrapCluster(config)
	if err := future.Error(); err != nil {
		log.Fatalf("Failed to bootstrap cluster: %v", err)
	}

	// Wait until the first node becomes the leader
	time.Sleep(5 * time.Second)

	// Have the other nodes join the cluster
	for i := 1; i < nodeCount; i++ {
		addVoterFuture := nodes[0].AddVoter(ids[i], transports[i].LocalAddr(), 0, 0)
		if err := addVoterFuture.Error(); err != nil {
			log.Fatalf("Failed to add voter for node %d: %v", i, err)
		}
	}

	log.Println("Raft cluster initialized and running")
	select {} // Keep the main goroutine running
}
