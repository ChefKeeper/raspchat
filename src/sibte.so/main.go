package main

/*
Copyright (c) 2015 Zohaib
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/natefinch/lumberjack.v2"

	"sibte.so/rasconfig"
	"sibte.so/rasweb"
	"sibte.so/rica"
    "sibte.so/rascluster"
)

func installSocketMux(mux *http.ServeMux, appConfig *rasconfig.ApplicationConfig, stateMachine rascluster.ClusterStateMachine) (err error) {
	err = nil
	s := rica.NewChatService(appConfig, stateMachine).WithRESTRoutes("/chat")

	mux.Handle("/chat", s)
    mux.Handle("/chat/", s)
	return
}

func installHTTPRoutes(mux *http.ServeMux, appConfig *rasconfig.ApplicationConfig, stateMachine rascluster.ClusterStateMachine) (err error) {
    routeHandlers := []rasweb.RouteHandler{
        rasweb.NewGifHandler(appConfig),
        rasweb.NewFileUploadHandler(appConfig),
        rasweb.NewClusterHandler(appConfig, stateMachine),
        rasweb.NewConfigRouteHandler(appConfig),
        rasweb.NewDirectPagesHandler(appConfig),
    }

	err = nil
	router := httprouter.New()

	// Register all routes
	for _, h := range routeHandlers {
		if err := h.Register(router); err != nil {
			log.Panic("Unable to register route")
		}
	}

	router.ServeFiles("/static/*filepath", http.Dir("./static"))
	mux.Handle("/", router)
	return
}

func parseArgs() (filePath string) {
	flag.StringVar(&filePath, "config", "", "Path to configuration file")
	flag.Parse()
	return
}

func main() {
    rica.RegisterMessageTypes()
	rasconfig.LoadApplicationConfig(parseArgs())
	conf := &rasconfig.CurrentAppConfig
	if conf.DBPath != "" {
		os.MkdirAll(conf.DBPath, os.ModePerm)
	}
	if conf.LogFilePath != "" {
		os.MkdirAll(conf.LogFilePath, os.ModePerm)
	}
	if conf.ClusterStatePath != "" {
		os.MkdirAll(conf.ClusterStatePath, os.ModePerm)
	}

	if conf.LogFilePath != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   conf.LogFilePath,
			MaxBackups: 3,
			MaxSize:    5,
			MaxAge:     15,
		})
	}

    stateMachine, err := rascluster.NewRaftStateMachine(
        conf.ClusterStatePath,
        conf.ClusterBindAddress,
        conf.ClusterPeers == nil || len(conf.ClusterPeers) < 1)

    if err != nil {
        log.Panic(err)
    }

	mux := http.NewServeMux()
	installSocketMux(mux, conf, stateMachine)
	installHTTPRoutes(mux, conf, stateMachine)
	server := &http.Server{
		Addr:    conf.BindAddress,
		Handler: mux,
	}

	log.Println("Starting server...", conf.BindAddress)
	log.Panic(server.ListenAndServe())
}
