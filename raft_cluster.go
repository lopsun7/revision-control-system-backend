package revision_control_system_backend

import (
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"
	"revision_control_system_backend/pb"
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

func main() {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize five Raft nodes
	nodes := make([]*server, 5)
	ports := []int{12000, 12001, 12002, 12003, 12004}
	for i := range nodes {
		nodeID := "node" + strconv.Itoa(i+1)
		ra, store := setupRaft(nodeID, ports[i])
		nodes[i] = &server{raft: ra, store: store}
		pb.RegisterRevisionControlServer(grpcServer, nodes[i])
	}

	log.Println("Server listening at", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
