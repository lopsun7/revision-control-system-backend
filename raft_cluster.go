package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

func setupRaft(nodeID string, raftPort int) (*raft.Raft, raft.Transport, raft.ServerID) {
	// Define base directory for this node
	baseDir := filepath.Join("tmp", "raft", nodeID)
	os.MkdirAll(baseDir, 0700)

	// Setup Raft configuration
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	// Setup network address and transport
	address := "127.0.0.1:" + strconv.Itoa(raftPort)
	addr, _ := net.ResolveTCPAddr("tcp", address)
	transport, _ := raft.NewTCPTransport(address, addr, 3, 10*time.Second, os.Stderr)

	// Setup log store, stable store, and snapshot store
	logStore, _ := raftboltdb.NewBoltStore(filepath.Join(baseDir, "raft.db"))
	stableStore, _ := raftboltdb.NewBoltStore(filepath.Join(baseDir, "stable.db"))
	snapshotStore, _ := raft.NewFileSnapshotStore(baseDir, 1, os.Stderr)

	// Create Raft node
	node, err := raft.NewRaft(config, nil, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatalf("Failed to create Raft node: %v", err)
	}

	return node, transport, raft.ServerID(nodeID)
}
