package rasweb

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"sibte.so/rascluster"
	"sibte.so/rasconfig"

	"github.com/julienschmidt/httprouter"
)

type clusterHandler struct {
	stateMachine rascluster.ClusterStateMachine
}

func NewClusterHandler() RouteHandler {
	return &clusterHandler{}
}

func (c *clusterHandler) Register(h *httprouter.Router) error {
	config := rasconfig.CurrentAppConfig
	sm, err := rascluster.NewRaftStateMachine(
		config.ClusterStatePath,
		config.ClusterBindAddress,
		config.ClusterPeers == nil || len(config.ClusterPeers) < 1)

	if err != nil {
		return err
	}

	c.stateMachine = sm
	h.POST("/cluster/add", c.addClusterPeer)
	h.GET("/cluster/join/:address", c.joinCluster)

	for _, peer := range config.ClusterPeers {
		log.Println("Requesting to join...", peer, c.requestJoinCluster(peer))
	}

	return nil
}

func (c *clusterHandler) addClusterPeer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !c.stateMachine.IsLeader() {
		w.WriteHeader(403)
		fmt.Fprintf(w, "Invalid state to perform operation, node is not a leader")
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid request %v", err)
		return
	}

	peer := r.Form.Get("peer")
	if failedPeers := c.stateMachine.AddPeers([]string{peer}); len(failedPeers) != 0 {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Unable to add peer %v because %v", peer, failedPeers[peer])
		return
	}

	fmt.Fprintf(w, "%v", peer)
}

func (c *clusterHandler) joinCluster(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if err := c.requestJoinCluster(params.ByName("peer")); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error %v", err)
		return
	}

	fmt.Fprintf(w, "OK")
}

func (c *clusterHandler) requestJoinCluster(peer string) error {
	peer_url := fmt.Sprintf("http://%s/cluster/add", peer)
	if _, err := http.PostForm(peer_url, url.Values{"peer": {c.stateMachine.Address()}}); err != nil {
		return err
	}

	return nil
}
