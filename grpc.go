package revision_control_system_backend

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/raft"
	pb "revision_control_system_backend/pb"
)

type server struct {
	raft  *raft.Raft
	store *FileStore
	pb.UnimplementedRevisionControlServer
}

func (s *server) CommitFile(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	cmd := Command{Filename: req.Filename, Content: req.Content}
	data, _ := json.Marshal(cmd)

	applyFuture := s.raft.Apply(data, 100)
	if err := applyFuture.Error(); err != nil {
		return &pb.CommitResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.CommitResponse{Success: true, Message: "Committed"}, nil
}

func (s *server) GetFile(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	versions, found := s.store.files[req.Filename]
	if !found || len(versions) == 0 {
		return &pb.GetResponse{Success: false, Message: "File not found"}, nil
	}
	latestVersion := versions[len(versions)-1] // 获取最新的版本
	return &pb.GetResponse{
		Content: latestVersion.Content, // 仅传递内容字符串
		Success: true,
		Message: "Success",
	}, nil
}
