package revision_control_system_backend

import (
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"
	"revision_control_system_backend/pb"
	"time"
)

func main() {
	ra, fileStore := setupRaft() // 修改 setupRaft 来返回 *FileStore

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRevisionControlServer(grpcServer, &server{raft: ra, store: fileStore})
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setupRaft() (*raft.Raft, *FileStore) {
	baseDir := filepath.Join("tmp", "raft")
	os.MkdirAll(baseDir, 0700)

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID("node1")

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:12000")
	transport, _ := raft.NewTCPTransport(addr.String(), addr, 3, 10*time.Second, os.Stderr)

	logStore, _ := raftboltdb.NewBoltStore(filepath.Join(baseDir, "raft.db"))
	stableStore, _ := raftboltdb.NewBoltStore(filepath.Join(baseDir, "stable.db"))
	snapshotStore, _ := raft.NewFileSnapshotStore(baseDir, 1, os.Stderr)

	fsm := NewFileStore()
	ra, _ := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	return ra, fsm
}
