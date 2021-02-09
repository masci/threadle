package logger

import (
	"log"

	"github.com/masci/threadle/intake"
)

var broker *intake.PubSub

// Init subscribes and processes the messages for the supported topics
func Init(b *intake.PubSub) {
	broker = b
	// log all the /api/v1/* endpoints
	for _, ep := range intake.GetV1Endpoints() {
		process(ep)
	}
	// log the /intake endpoint
	process(intake.IntakeEndpointV1)
}

// process reads the message from the broker and logs to stderr
func process(topic string) {
	go func() {
		for msg := range broker.Subscribe(topic) {
			log.Println(topic)
			log.Println(string(msg))
		}
	}()
}
