package logger

import (
	"time"

	"github.com/masci/threadle/intake"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	broker *intake.PubSub
	logger zerolog.Logger
)

// Plugin implements plugins.Plugin
type Plugin struct{}

// Start subscribes and processes the messages for the supported topics
func (*Plugin) Start(b *intake.PubSub) {
	broker = b

	// Setup the logger
	if viper.GetBool("plugins.logger.ecs_compatible") {
		zerolog.TimeFieldFormat = time.RFC3339
		zerolog.TimestampFieldName = "@timestamp"
		zerolog.MessageFieldName = "message"
		zerolog.LevelFieldName = "log.level"
		log.Logger = log.With().Str("ecs.version", "1.6.0").Logger()
	}

	// Log all the /api/v1/* endpoints
	for _, ep := range intake.GetV1Endpoints() {
		process(ep)
	}
	// Log the /intake endpoint
	process(intake.IntakeEndpointV1)
}

// process reads the message from the broker and logs to stderr
func process(topic string) {
	go func() {
		for msg := range broker.Subscribe(topic) {
			log.Info().Str("topic", topic).RawJSON("message", msg).Msg("")
		}
	}()
}
