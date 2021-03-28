package main

import (
	"os"

	"github.com/masci/threadle/intake"
	"github.com/masci/threadle/output"
	"github.com/masci/threadle/plugins"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	// output plugins
	"github.com/masci/threadle/plugins/elasticsearch"
	"github.com/masci/threadle/plugins/logger"
)

func main() {
	// Define and parse command args
	verbosity := pflag.IntP("verbose", "v", 1, "set verbosity level: 0 silent, 1 normal, 2 debug")
	help := pflag.BoolP("help", "h", false, "print args help")
	pflag.Parse()

	// Print the help message and exit if --help is passed
	if *help == true {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	// Configure cmdline output facilities
	output.Init(*verbosity)

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
			output.INFO.Println("Initializing plugin:", k)
			p.Start(intake.MsgBroker)
		}
	}

	// Start the HTTP server, block until shutdown
	intake.Serve()
	os.Exit(0)
}

func initConfig() {
	viper.SetDefault("port", "3060")

	viper.SetConfigName("threadle.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		output.FATAL.Fatalf("Fatal error: %s", err)
	}
	output.DEBUG.Println("Config file found at", viper.ConfigFileUsed())
}
