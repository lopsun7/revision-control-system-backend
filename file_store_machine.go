package revision_control_system_backend

import (
	"encoding/json"
	"github.com/hashicorp/raft"
	"io"
)

type Command struct {
	Filename  string
	Author    string
	Content   string
	Timestamp int64
}
type FileVersion struct {
	Content   string
	Timestamp int64
}

type FileStore struct {
	files map[string][]FileVersion
}

func NewFileStore() *FileStore {
	return &FileStore{
		files: make(map[string][]FileVersion),
	}
}

func (f *FileStore) Apply(log *raft.Log) interface{} {
	var cmd Command
	err := json.Unmarshal(log.Data, &cmd)
	if err != nil {
		return nil
	}
	versions, found := f.files[cmd.Filename]
	if !found {
		versions = []FileVersion{}
	}
	newVersion := FileVersion{Content: cmd.Content, Timestamp: cmd.Timestamp}
	versions = append(versions, newVersion)
	f.files[cmd.Filename] = versions
	return nil
}

func (f *FileStore) Snapshot() (raft.FSMSnapshot, error) {
	return &FileSnapshot{store: f}, nil
}

func (f *FileStore) Restore(rc io.ReadCloser) error {
	// Assume the restore logic is implemented here
	return nil
}

type FileSnapshot struct {
	store *FileStore
}

func (fs *FileSnapshot) Persist(sink raft.SnapshotSink) error {
	// Implement the persistence logic here
	return sink.Close()
}

func (fs *FileSnapshot) Release() {
	// Any cleanup can be performed here
}
