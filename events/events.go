package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/troneras/gorews/data"
	"github.com/troneras/gorews/events/config"
	"github.com/troneras/gorews/logger"
	log "github.com/troneras/gorews/logger"
)

type ExternalServerMessage struct {
	EventType string
	Data      *data.Message
}

// buffer size for the channel (needs to handle all open/close websockets events of a high traffic website)
var bufferSize = 1000000 // memory usage: 1000000 * 16 bytes = 16 MB

var externalServerChan = make(chan ExternalServerMessage, bufferSize)

var conf *config.Config

func Configure() {
	log.Info("[EVENTS] Configuring events")
	conf = config.Configure()
}

func SendEvent(eventType string, data *data.Message) {
	externalServerChan <- ExternalServerMessage{
		EventType: eventType,
		Data:      data,
	}
}

func HandleExternalServerMessages() {
	for msg := range externalServerChan {
		log.Debug("[EVENTS] Received event ", log.Fields{"event": msg.EventType, "data": msg.Data})
		for _, event := range msg.Data.EventsSubscribed {
			if event == msg.EventType {
				sendToExternalServer(msg.EventType, msg.Data)
			}
		}
	}
}

func sendToExternalServer(eventType string, data *data.Message) {
	// Logic to send a request to the external server
	// Use an HTTP client to send data
	log.Debug("[EVENTS] Sending event " + eventType + " to external server")

	// Example of how to post json data to an external servcr with basic auth
	url := fmt.Sprintf("https://%s/phive/modules/Licensed/ajax.php", data.Domain)

	// body := map[string]string{"key": "value"}
	body := map[string]string{
		"return_format":   "json",
		"domain":          data.Domain,
		"websocket_event": data.Channel,
		"event":           eventType,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Error("[EVENTS] Error marshalling json body", log.Fields{"error": err, "body": body})
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error("[EVENTS] Error creating request", log.Fields{"error": err, "url": url})
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(jsonBody)))
	req.SetBasicAuth(conf.BasicUser, conf.BasicPass)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("[EVENTS] Error sending request", log.Fields{"error": err, "url": url})
		return
	}
	defer resp.Body.Close()
	log.Debug("[EVENTS] Response Status:", log.Fields{"status": resp.Status})
	log.Debug("[EVENTS] Response Headers:", log.Fields{"headers": resp.Header})

	// log the return body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("[EVENTS] Error reading response body", log.Fields{"error": err, "url": url, "body": body})
		return
	}
	logger.Debug("[EVENTS] Response Body:", log.Fields{"body": string(bodyBytes)})
}
