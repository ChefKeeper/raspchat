package rascluster

import "time"

type UpdateStateListener interface {
    OnStateUpdated([]byte) error
}

type ClusterStateMachine interface {
	AddPeers([]string) map[string]error
	Ping() error
	ApplyMessage(msg []byte, duration time.Duration) error

    OnUpdate(UpdateStateListener)
    OffUpdate(UpdateStateListener)

	IsLeader() bool
	Address() string
	Leader() string
}
