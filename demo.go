package revision_control_system_backend

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
)

// 定义文件版本信息
type FileVersion struct {
	Content string
	Author  string
	Time    time.Time
}

// 文件存储结构
type FileStore struct {
	Files map[string][]FileVersion
	raft  *raft.Raft
}

func NewFileStore() *FileStore {
	// Raft 配置部分
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID("node1")

	// 创建用于存储和快照的目录
	os.MkdirAll("./raftdb", os.ModePerm)
	logStore, err := raftboltdb.NewBoltStore("./raftdb/raft.db")
	if err != nil {
		log.Fatalf("failed to create log store: %s", err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore("./raftdb", 1, nil)
	if err != nil {
		log.Fatalf("failed to create snapshot store: %s", err)
	}

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Fatalf("failed to resolve TCP address: %s", err)
	}

	transport, err := raft.NewTCPTransport(addr.String(), addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalf("failed to create transport: %s", err)
	}

	ra, err := raft.NewRaft(config, nil, logStore, logStore, snapshotStore, transport)
	if err != nil {
		log.Fatalf("failed to create raft: %s", err)
	}

	return &FileStore{
		Files: make(map[string][]FileVersion),
		raft:  ra,
	}
}

func main() {
	fs := NewFileStore()

	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// 此处注册 gRPC 服务，略

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
