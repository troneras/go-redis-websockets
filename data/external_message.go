package data

type ExternalServerMessage struct {
	EventType string
	Data      *Message
}

var ExternalServerChan = make(chan ExternalServerMessage)

func SendEvent(eventType string, data *Message) {
	ExternalServerChan <- ExternalServerMessage{
		EventType: eventType,
		Data:      data,
	}
}
