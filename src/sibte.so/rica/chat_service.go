package rica

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"sibte.so/rasconfig"
    "sibte.so/rascluster"
)

type ChatService struct {
    sync.Mutex
    groupInfo           GroupInfoManager
    clusterStateMachine rascluster.ClusterStateMachine
    chatStore           ChatLogStore
    nickRegistry        *NickRegistry
    wsUpgrader          *websocket.Upgrader
    gcmWorker           *GCMWorker
    httpMux             *http.ServeMux
}

func NewChatService(appConfig *rasconfig.ApplicationConfig, stateMachine rascluster.ClusterStateMachine) *ChatService {
	initChatHandlerTypes()
	store, e := NewChatLogStore(rasconfig.CurrentAppConfig.DBPath+"/chats.bolt.db", []byte("chats"))
	allowedOrigins := appConfig.AllowedOrigins

	if e != nil {
		log.Panic(e)
	}

	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	if len(allowedOrigins) > 0 {
		wsUpgrader.CheckOrigin = func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" || allowedOrigins == nil || len(allowedOrigins) == 0 {
				return true
			}

			for _, item := range allowedOrigins {
				if strings.Compare(item, origin) == 0 {
					return true
				}
			}

			log.Println("Denying connection due to missing origin " + origin)
			return false
		}
	}

	return &ChatService{
		groupInfo:              NewInMemoryGroupInfo(),
		nickRegistry:           NewNickRegistry(),
        clusterStateMachine:    stateMachine,
		gcmWorker:              NewGCMWorker(rasconfig.CurrentAppConfig.GCMToken),
		chatStore:              NewReplicatedChatLogStore(store, stateMachine),
		wsUpgrader:             wsUpgrader,
	}
}

func (c *ChatService) WithRESTRoutes(prefix string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(prefix+"/api/", c.httpRoutes(prefix+"/api", httprouter.New()))
	c.httpMux = mux
	return c
}

func (c *ChatService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/chat/api") {
		c.httpMux.ServeHTTP(w, req)
		return
	}

	c.upgradeConnectionToWebSocket(w, req)
}

func (c *ChatService) httpRoutes(prefix string, router *httprouter.Router) http.Handler {
	router.POST(prefix+"/push", c.onPushPost)
	router.POST(prefix+"/register", c.onPushSubscribe)

	router.GET(prefix+"/channel/:id/message", c.onGetChatHistory)
	router.GET(prefix+"/channel/:id/message/:msg_id", c.onGetChatMessage)
	router.GET(prefix+"/channel", c.onGetChannels)
	return router
}

func (c *ChatService) upgradeConnectionToWebSocket(w http.ResponseWriter, req *http.Request) bool {
	conn, err := c.wsUpgrader.Upgrade(w, req, nil)
	if err == nil {
		transporter := NewWebsocketMessageTransport(conn)
		handler := NewChatHandler(c.nickRegistry, c.groupInfo, transporter, c.chatStore)
		go handler.Loop()
		return true
	}

	log.Println("Error upgrading connection...", err)
	return false
}

func (c *ChatService) onPushSubscribe(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	token := req.FormValue("gcm_sub_token")
	if token == "" {
		fmt.Fprint(w, "false")
		return
	}

	transporter := NewGCMTransport(token, c.gcmWorker)
	handler := NewChatHandler(c.nickRegistry, c.groupInfo, transporter, c.chatStore)
	go handler.Loop()
	fmt.Fprint(w, "true")
}

func (c *ChatService) onPushPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	token := req.FormValue("gcm_sub_token")
	t := NewGCMTransport(token, c.gcmWorker)
	if msg, err := ioutil.ReadAll(req.Body); req.Method == "POST" && err == nil {
		t.PostMessage(string(msg))
		fmt.Fprint(w, "true")
		return
	}

	fmt.Fprint(w, "false")
}

func (c *ChatService) onGetChatHistory(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	groupID := p.ByName("id")

	queryParams := req.URL.Query()
	var offset uint = 0
	var limit uint = 20
	startID := queryParams.Get("start_id")

	if o, err := strconv.ParseUint(queryParams.Get("offset"), 10, 32); err == nil {
		offset = uint(o)
	}

	if l, err := strconv.ParseUint(queryParams.Get("limit"), 10, 32); err == nil {
		limit = uint(l)
	}

	chatLog, err := c.chatStore.GetMessagesFor(groupID, startID, offset, limit)
	if err == nil {
		response := make(map[string]interface{})
		response["limit"] = limit
		response["offset"] = offset
		response["messages"] = chatLog
		response["start_id"] = startID
		response["id"] = groupID
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMessage{
			Error: err.Error(),
		})
	}
}

func (c *ChatService) onGetChatMessage(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
}

func (c *ChatService) onGetChannels(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
}
