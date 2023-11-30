package main

import (
	gohttp "net/http"
	"os"
	"os/signal"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/troneras/gorews/api"
	cfgapi "github.com/troneras/gorews/api/config"
	"github.com/troneras/gorews/events"
	"github.com/troneras/gorews/http"
	log "github.com/troneras/gorews/logger"
	"github.com/troneras/gorews/redis"
	"github.com/troneras/gorews/websockets"
)

var apiconf *cfgapi.Config

func configure() {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		godotenv.Load(".env.production")
	} else {
		godotenv.Load(".env.development")
	}

	log.Configure()
	apiconf = cfgapi.Configure()
	redis.Configure()
	websockets.Configure()
	http.Configure()
	events.Configure()
}

func main() {
	configure()

	exitCh := make(chan int)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	cb := func(r gohttp.Handler) {
		api.CreateApi(apiconf, r.(*pat.Router))
	}

	go http.Listen(apiconf.APIBindAddr, exitCh, cb)
	go events.HandleExternalServerMessages()

	defer websockets.CloseAllHubs()

	select {
	case <-exitCh:
		log.Info("[MAIN] Exiting")
	case <-signalCh:
		log.Info("[MAIN] Interrupt signal received, exiting")
	}
	os.Exit(0)
}

/*

Add some random content to the end of this file, hopefully tricking GitHub
into recognising this as a Go repo instead of Makefile.

A gopher, ASCII art style - borrowed from
https://gist.github.com/belbomemo/b5e7dad10fa567a5fe8a

          ,_---~~~~~----._
   _,,_,*^____      _____``*g*\"*,
  / __/ /'     ^.  /      \ ^@q   f
 [  @f | @))    |  | @))   l  0 _/
  \`/   \~____ / __ \_____/    \
   |           _l__l_           I
   }          [______]           I
   ]            | | |            |
   ]             ~ ~             |
   |                            |
    |                           |

*/
