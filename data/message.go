package data

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"

	"github.com/troneras/gorews/logger"
)

// Message represents a parsed SMTP message
type Message struct {
	Path             string
	Channel          string
	Domain           string
	EventsSubscribed []string
	Handler          string
}

// NewMessageFromRequest creates a new Message from the request uri and params  /id/:id/[tag/:tag]][?events=open,close&handler=handlerName]]]
func NewMessageFromURL(sha1_secret string, r *http.Request) *Message {
	// the channel is the sha1 of the /id/:id/[tag/:tag] part
	url := r.URL
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s%s", sha1_secret, url.Path)))
	channel := fmt.Sprintf("%x", h.Sum(nil))
	// ev is a string like "open,close", split it into a slice of strings
	events := strings.Split(url.Query().Get("events"), ",")
	// the handler is the redis key that stores the events callback function on the backend
	domain := getOriginFromRequest(r)
	handler := url.Query().Get("handler")

	return &Message{
		Path:             url.Path,
		Channel:          channel,
		EventsSubscribed: events,
		Domain:           domain,
		Handler:          handler,
	}
}

func getOriginFromRequest(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin == "" {
		logger.Println("[APIv1] Origin header not set")
		origin = fmt.Sprintf("http://%s", r.Host)
	}
	origin = strings.Split(origin, "//")[1]
	return origin
}

// String returns a string representation of the Message
func (m *Message) String() string {
	events := strings.Join(m.EventsSubscribed, ", ")
	return "Path: " + m.Path + ", Channel: " + m.Channel + ", Events: " + events + ", Handler: " + m.Handler + ", Domain: " + m.Domain
}
