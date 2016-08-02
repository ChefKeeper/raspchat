package rascluster

type ClusterStateMachine interface {
	AddPeers([]string) map[string]error
	Ping() error
	IsLeader() bool
	Address() string
	Leader() string
}
