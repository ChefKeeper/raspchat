package rascluster

type ClusterStateMachine interface {
	AddPeers([]string) map[string]error
	IsLeader() bool
	Address() string
}
