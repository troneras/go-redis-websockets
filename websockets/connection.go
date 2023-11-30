package websockets

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/troneras/gorews/data"
	log "github.com/troneras/gorews/logger"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer. Set to minimum allowed value as we don't expect the client to send non-control messages.
	maxMessageSize = 256
)

type connection struct {
	hub  *Hub
	ws   *websocket.Conn
	send chan interface{}
	info *data.Message
}

func (c *connection) readLoop() {
	defer func() {
		c.hub.unregisterChan <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, data, err := c.ws.NextReader()
		if err != nil {
			log.Error("[WEBSOCKET] Error reading from connection", log.Fields{"error": err})
			return
		}
		switch mt {
		case websocket.TextMessage:
			c.hub.messages <- data
		case websocket.CloseMessage:
			return
		default:
			log.Warn("[WEBSOCKET] Received unexpected message type", log.Fields{"type": mt})
		}
	}
}

func (c *connection) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.writeControl(websocket.CloseMessage)
				return
			}
			if err := c.writeText(message.(string)); err != nil {
				log.Error("[WEBSOCKET] Error writing message", log.Fields{"error": err, "message": message})
				return
			}
		case <-ticker.C:
			if err := c.writeControl(websocket.PingMessage); err != nil {
				log.Error("[WEBSOCKET] Error writing ping message", log.Fields{"error": err})
				return
			}
		}
	}
}

func (c *connection) Close() {
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	c.ws.WriteMessage(websocket.CloseMessage, message)
	c.ws.Close()
	c.hub.unregisterChan <- c
}
func (c *connection) writeText(message string) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *connection) writeControl(messageType int) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, []byte{})
}
