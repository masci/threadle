package elasticsearch

import (
	"regexp"

	"github.com/masci/threadle/intake"
)

type filters []*regexp.Regexp

// Init subscribes and processes the messages for the supported topics
func Init(b *intake.PubSub) {
	// configure filters
	exclude := filters{
		regexp.MustCompile(`system\..`),         // system.*
		regexp.MustCompile(`datadog\.agent\..`), // datadog.agent.*
	}
	// subcsribe to messages
	readV1Metrics(b, exclude)
}
