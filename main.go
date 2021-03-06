package main

import (
	"fmt"
	"os"

	"github.com/masci/threadle/intake"
	"github.com/masci/threadle/plugins"
	"github.com/spf13/viper"

	// output plugins
	"github.com/masci/threadle/plugins/elasticsearch"
	"github.com/masci/threadle/plugins/logger"
)

func main() {
	// Bootstrap config, this has to be called first
	initConfig()

	// Define the available output plugins
	plugins := map[string]plugins.Plugin{
		"logger":        &logger.Plugin{},
		"elasticsearch": &elasticsearch.Plugin{},
	}

	// Load the configured output plugins
	for k := range viper.GetStringMap("plugins") {
		if p, found := plugins[k]; found {
			fmt.Println("Initializing plugin:", k)
			p.Start(intake.MsgBroker)
		}
	}

	// Start the HTTP server, block until shutdown
	intake.Serve()
	os.Exit(0)
}

func initConfig() {
	viper.SetDefault("port", "8080")

	viper.SetConfigName("threadle.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error: %s", err))
	}
}
