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
    appConfig    *rasconfig.ApplicationConfig
}

func NewClusterHandler(config *rasconfig.ApplicationConfig, stateMachine rascluster.ClusterStateMachine) RouteHandler {
    return &clusterHandler{
        appConfig: config,
        stateMachine: stateMachine,
    }
}

func (c *clusterHandler) Register(h *httprouter.Router) error {
    h.POST("/cluster/add", c.addClusterPeer)
    h.GET("/cluster/join/:address", c.joinCluster)
    h.GET("/cluster/ping", c.pingCluster)

    for _, peer := range c.appConfig.ClusterPeers {
        log.Println("Requesting to join...", peer, c.requestJoinCluster(peer))
    }

    return nil
}

func (c *clusterHandler) addClusterPeer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    if !c.stateMachine.IsLeader() {
        w.WriteHeader(403)
        fmt.Fprint(w, "Invalid state to perform operation, node is not a leader")
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

    fmt.Fprint(w, "OK")
}

func (c *clusterHandler) pingCluster(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    if !c.stateMachine.IsLeader() {
        w.WriteHeader(409)
        fmt.Fprintf(w, "Not a leader current leader is %v", c.stateMachine.Leader())
        return
    }

    if err := c.stateMachine.Ping(); err != nil {
        w.WriteHeader(500)
        fmt.Fprintf(w, "Error pinging %v", err)
        return
    }

    fmt.Fprint(w, "Success")
}

func (c *clusterHandler) requestJoinCluster(peer string) error {
    peer_url := fmt.Sprintf("http://%s/cluster/add", peer)
    if _, err := http.PostForm(peer_url, url.Values{"peer": {c.stateMachine.Address()}}); err != nil {
        return err
    }

    return nil
}
