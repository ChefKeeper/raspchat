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
    callbacks []UpdateStateListener
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
        callbacks: make([]UpdateStateListener, 0),
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

func (s *raftStateMachine) Leader() string {
	return s.raft.Leader()
}

func (s *raftStateMachine) ApplyMessage(msg []byte, duration time.Duration) error {
    if f := s.raft.VerifyLeader(); f != nil && f.Error() != nil {
        return f.Error()
    }

    if f := s.raft.Apply(msg, duration); f != nil && f.Error() != nil {
        return f.Error()
    }

    return nil
}

func (s *raftStateMachine) Ping() error {
    return s.ApplyMessage([]byte("PING"), 1*time.Second)
}

func (s *raftStateMachine) OnUpdate(callback UpdateStateListener) {
    s.callbacks = append(s.callbacks, callback)
}

func (s *raftStateMachine) OffUpdate(callback UpdateStateListener) {
    found := -1
    for i, cb := range s.callbacks {
        if callback == cb {
            found = i
        }
    }

    if found != -1 {
        s.callbacks = append(s.callbacks[:found], s.callbacks[found+1:]...)
    }
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (s *raftStateMachine) Apply(logEntry *raft.Log) interface{} {
	log.Println("Applying log type:", logEntry.Type, "index:", logEntry.Index, "term:", logEntry.Term)

    // Only play update state in case of command log
    if logEntry.Type == raft.LogCommand {
        log.Println("About to invoke", len(s.callbacks), "subscriptions")
        for _, cb := range s.callbacks {
            if err := cb.OnStateUpdated(logEntry.Data); err != nil {
                log.Println("ERROR =", err)
            }
        }
    }

	return nil
}

// Snapshot is used to support log compaction. This call should
// return an FSMSnapshot which can be used to save a point-in-time
// snapshot of the FSM. Apply and Snapshot are not called in multiple
// threads, but Apply will be called concurrently with Persist. This means
// the FSM should be implemented in a fashion that allows for concurrent
// updates while a snapshot is happening.
func (s *raftStateMachine) Snapshot() (raft.FSMSnapshot, error) {
	log.Println("Snapshot")
	return nil, nil
}

// Restore is used to restore an FSM from a snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
func (s *raftStateMachine) Restore(reader io.ReadCloser) error {
	log.Println("Restore")
	return nil
}
