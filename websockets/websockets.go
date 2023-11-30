package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/troneras/gorews/data"
	"github.com/troneras/gorews/events"
	log "github.com/troneras/gorews/logger"
	"github.com/troneras/gorews/redis"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	redisChannel   string
	rr             *redis.RedisReader
	rw             *redis.RedisWriter
	upgrader       *websocket.Upgrader
	connections    map[*connection]bool
	messages       chan interface{}
	registerChan   chan *connection
	unregisterChan chan *connection
}

var hubs map[string]*Hub

var exitCh chan string

func Configure() {
	hubs = make(map[string]*Hub)
	exitCh = make(chan string)

	go func() {
		for ch := <-exitCh; ch != ""; ch = <-exitCh {
			log.Debug("[WEBSOCKET] Exiting for channel", log.Fields{"channel": ch})
			delete(hubs, ch)
		}
	}()
}

func GetHubForChannel(channel string) *Hub {
	if hub, ok := hubs[channel]; ok {
		return hub
	}
	hubs[channel] = NewHub(channel)
	return hubs[channel]
}

func NewHub(redisChannel string) *Hub {
	hub := &Hub{
		upgrader:       &upgrader,
		redisChannel:   redisChannel,
		rr:             redis.NewRedisReader(redisChannel),
		rw:             redis.NewRedisWriter(redisChannel),
		connections:    make(map[*connection]bool),
		messages:       make(chan interface{}),
		registerChan:   make(chan *connection),
		unregisterChan: make(chan *connection),
	}
	go hub.run()
	return hub
}

func (h *Hub) run() {
	defer func() {
		log.Debug("[WEBSOCKET] Closing hub for channel", log.Fields{"channel": h.redisChannel})
		for c := range h.connections {
			h.unregister(c)
		}
	}()
	for {
		select {
		case c := <-h.registerChan:
			h.register(c)
		case c := <-h.unregisterChan:
			h.unregister(c)
		case m, ok := <-h.rr.MessageChan:
			if !ok {
				log.Debug("[WEBSOCKET] Redis reader channel closed", log.Fields{"channel": h.redisChannel})
				return
			}
			// when a message is received from redis, send it to all connections
			log.Debug("[WEBSOCKET] sending message to all connections", log.Fields{"channel": h.redisChannel, "message": m})
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					h.unregister(c)
				}
			}
		case m := <-h.messages:
			log.Debug("[WEBSOCKET] Received message from API", log.Fields{"channel": h.redisChannel, "message": m})
			// TODO do we really want to broadcast messages from the API to all connections without filtering?
			// h.rr.MessageChan <- m
			/* for c := range h.connections {
				select {
				case c.send <- m:
				default:
					h.unregister(c)
				}
			} */
		}
	}
}

func (h *Hub) register(c *connection) {
	log.Debug("[WEBSOCKET] Registering connection")
	h.connections[c] = true
	events.SendEvent("open", c.info)
}

func (h *Hub) unregister(c *connection) {
	if _, ok := h.connections[c]; ok {
		events.SendEvent("close", c.info)
		close(c.send)
		delete(h.connections, c)
		log.Debug("[WEBSOCKET] Unregistered connection", log.Fields{"channel": h.redisChannel, "remaining": len(h.connections)})
		// if there are no more connections, close the redis reader and writer and exit
		if len(h.connections) == 0 {
			log.Debug("[WEBSOCKET] No more connections, closing redis reader and writer")
			h.rr.Close()
			h.rw.Close()
			log.Debug("[WEBSOCKET] Sending exit signal to hub ", log.Fields{"channel": h.redisChannel})
			exitCh <- h.redisChannel
		}
	} else {
		log.Warn("[WEBSOCKET] Connection not found when trying to unregister")
	}
}

func (h *Hub) Serve(w http.ResponseWriter, r *http.Request, info *data.Message) {
	log.Debug("[WEBSOCKET] Upgrading connection")
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("[WEBSOCKET] Error upgrading connection", log.Fields{"error": err})
		return
	}
	log.Debug("[WEBSOCKET] Upgraded connection")

	c := &connection{
		hub:  h,
		ws:   ws,
		send: make(chan interface{}, 256),
		info: info,
	}

	log.Debug("[WEBSOCKET] Registering connection")

	h.registerChan <- c

	log.Debug("[WEBSOCKET] Serving connection")

	go c.writeLoop()
	go c.readLoop()
}

func CloseAllHubs() {
	log.Debug("[WEBSOCKET] Closing all hubs")
	for _, hub := range hubs {
		for c := range hub.connections {
			hub.unregister(c)
		}
	}
	log.Debug("[WEBSOCKET] Closed all hubs")
}
