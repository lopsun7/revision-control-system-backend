package main

import (
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func setupRaft(nodeID string, raftPort int) (*raft.Raft, *FileStore) {
	baseDir := filepath.Join("tmp", "raft", nodeID)
	os.MkdirAll(baseDir, 0700)

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	address := "127.0.0.1:" + strconv.Itoa(raftPort)
	addr, _ := net.ResolveTCPAddr("tcp", address)
	transport, _ := raft.NewTCPTransport(address, addr, 3, 10*time.Second, os.Stderr)

	logStore, _ := raftboltdb.NewBoltStore(filepath.Join(baseDir, "raft.db"))
	stableStore, _ := raftboltdb.NewBoltStore(filepath.Join(baseDir, "stable.db"))
	snapshotStore, _ := raft.NewFileSnapshotStore(baseDir, 1, os.Stderr)

	fsm := NewFileStore()
	ra, _ := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	return ra, fsm
}
