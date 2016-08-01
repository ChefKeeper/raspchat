package rascluster

import (
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

type raftStateMachine struct {
	directory string
	bind      *net.TCPAddr
	raft      *raft.Raft
}

func NewRaftStateMachine(directory, bind string, leader bool) (ClusterStateMachine, error) {
	log.Println("Starting RaftStateMachine with", directory, bind, leader)

	config := raft.DefaultConfig()
	addr, err := net.ResolveTCPAddr("tcp", bind)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(bind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	peersStore := raft.NewJSONPeers(directory, transport)
	if leader {
		config.EnableSingleNode = true
		config.DisableBootstrapAfterElect = true
	}

	snapshotStore, err := raft.NewFileSnapshotStore(directory, 2, os.Stderr)
	if err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(path.Join(directory, "raft.log.bdb"))
	if err != nil {
		return nil, err
	}

	sm := &raftStateMachine{
		directory: directory,
		bind:      addr,
		raft:      nil,
	}

	sm.raft, err = raft.NewRaft(config, sm, logStore, logStore, snapshotStore, peersStore, transport)
	if err != nil {
		return nil, err
	}

	return sm, nil
}

func (s *raftStateMachine) AddPeers(peers []string) map[string]error {
	errorsMap := make(map[string]error)
	for _, p := range peers {
		log.Println("Joining peer...", p)
		f := s.raft.AddPeer(p)
		if f.Error() != nil {
			errorsMap[p] = f.Error()
		}
	}

	return errorsMap
}

func (s *raftStateMachine) IsLeader() bool {
	return s.raft.State() == raft.Leader
}

func (s *raftStateMachine) Address() string {
	return s.bind.String()
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (s *raftStateMachine) Apply(log *raft.Log) interface{} {
	return nil
}

// Snapshot is used to support log compaction. This call should
// return an FSMSnapshot which can be used to save a point-in-time
// snapshot of the FSM. Apply and Snapshot are not called in multiple
// threads, but Apply will be called concurrently with Persist. This means
// the FSM should be implemented in a fashion that allows for concurrent
// updates while a snapshot is happening.
func (s *raftStateMachine) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

// Restore is used to restore an FSM from a snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
func (s *raftStateMachine) Restore(reader io.ReadCloser) error {
	return nil
}
