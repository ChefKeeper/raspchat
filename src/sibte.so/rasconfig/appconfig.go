package rasconfig

import (
	"io/ioutil"
	"log"
	"path"

	"encoding/json"
)

type ApplicationConfig struct {
	BindAddress        string            `json:"bind_address"`
	LogFilePath        string            `json:"log_file"`
	DBPath             string            `json:"db_path"`
	AllowHotRestart    bool              `json:"allow_hot_reboot"`
	GCMToken           string            `json:"gcm_token"`
	AllowedOrigins     []string          `json:"allowed_origins"`
	ExternalSignIn     map[string]string `json:"external_sign_in"`
	WebSocketURL       string            `json:"websocket_url"`
	WebSocketSecureURL string            `json:"websocketsecure_url"`
	HasAuthProviders   bool              `json:"has_auth_providers"`
	UploaderConfig     map[string]string `json:"uploader_config"`

	ClusterBindAddress string   `json:"cluster_bind_address"`
	ClusterPeers       []string `json:"cluster_peers"`
	ClusterStatePath   string   `json:"cluster_state_path"`
}

var CurrentAppConfig ApplicationConfig

func LoadApplicationConfig(filePath string) {
	dir, err := ioutil.TempDir("", "raspchat")
	if err != nil {
		log.Fatal(err)
	}

	conf := &CurrentAppConfig
	if filePath == "" {
		conf.AllowHotRestart = false
		conf.BindAddress = ":8080"
		conf.DBPath = dir
		conf.LogFilePath = ""
		conf.AllowedOrigins = make([]string, 0)
		conf.ExternalSignIn = make(map[string]string)
		conf.HasAuthProviders = false
		conf.WebSocketURL = ""
		conf.WebSocketSecureURL = ""

		conf.UploaderConfig = make(map[string]string)
		conf.UploaderConfig["provider"] = "local"
		conf.UploaderConfig["disk_storage_path"] = dir

		conf.ClusterPeers = make([]string, 0)
		conf.ClusterBindAddress = ":5000"
		conf.ClusterStatePath = path.Join(dir, "statemachine")
		return
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(content, &CurrentAppConfig); err != nil {
		log.Panic(err)
	}

	conf.HasAuthProviders = len(conf.ExternalSignIn) != 0
	log.Println("=== Loaded configuration")
	log.Println(CurrentAppConfig)
}
