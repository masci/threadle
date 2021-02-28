package plugins

import "github.com/masci/threadle/intake"

// Plugin is the interface provided by any plugin
type Plugin interface {
	Init(*intake.PubSub)
}
