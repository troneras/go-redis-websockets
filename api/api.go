package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/pat"
	"github.com/troneras/gorews/api/config"
	"github.com/troneras/gorews/data"
	log "github.com/troneras/gorews/logger"
	"github.com/troneras/gorews/websockets"
)

type APIv1 struct {
	config      *config.Config
	messageChan chan string
}

func CreateApi(conf *config.Config, r *pat.Router) *APIv1 {
	apiv1 := &APIv1{
		config:      conf,
		messageChan: make(chan string),
	}
	log.Info("[APIv1] Creating APIv1", log.Fields{"path": conf.ApiPath})
	r.Path(conf.ApiPath + "/id/{id}/tag/{tag}/").Methods("GET").HandlerFunc(apiv1.websocket)
	r.Path(conf.ApiPath + "/tag/{tag}/").Methods("GET").HandlerFunc(apiv1.websocket)

	return apiv1
}

func (apiv1 *APIv1) websocket(w http.ResponseWriter, r *http.Request) {
	log.Debug("[APIv1] /websocket called")
	// remove the conf.ApiPath from the url
	r.URL.Path = strings.Replace(r.URL.Path, apiv1.config.ApiPath, "", 1)

	msg := data.NewMessageFromURL(apiv1.config.Sha1Secret, r)
	log.Debug("[API] Received message", log.Fields{"function": "websocket", "msg": msg})

	channel := msg.Channel

	specificHub := websockets.GetHubForChannel(channel)

	specificHub.Serve(w, r, msg)
}
