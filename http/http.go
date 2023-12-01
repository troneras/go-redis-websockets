package http

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/troneras/gorews/http/config"
	log "github.com/troneras/gorews/logger"
)

var conf *config.Config

func Configure() {
	conf = config.Configure()
}
func SetHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("[HTTP] ****** NEW REQUEST ******", log.Fields{"method": r.Method, "url": r.URL.Path})
		h.ServeHTTP(w, r)
	})

}

// Listen binds to httpBindAddr
func Listen(httpBindAddr string, exitCh chan int, registerCallback func(http.Handler)) {

	log.Debug("[HTTP] Binding to address", log.Fields{"addr": httpBindAddr})

	pat := pat.New()
	registerCallback(pat)

	handler := SetHandler(pat)

	var err error

	if conf.UseTLS {
		log.Debug("[HTTP] Using TLS", log.Fields{"cert": conf.SSLCert, "key": conf.SSLKey})
		err = http.ListenAndServeTLS(httpBindAddr, conf.SSLCert, conf.SSLKey, handler)
	} else {
		err = http.ListenAndServe(httpBindAddr, handler)
	}
	if err != nil {
		log.Error("[HTTP] Error binding to address", log.Fields{"addr": httpBindAddr, "error": err})
	}
}
