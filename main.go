package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/masci/threadle/intake"
	"github.com/masci/threadle/plugins"
	"github.com/spf13/viper"

	// output plugins
	"github.com/masci/threadle/plugins/elasticsearch"
)

func main() {
	// bootstrap config
	initConfig()

	plugins := map[string]plugins.Plugin{
		"logger":        &elasticsearch.Plugin{},
		"elasticsearch": &elasticsearch.Plugin{},
	}

	// Initialize the intake
	intake.Init()
	// Configure the output plugins
	for k := range viper.GetStringMap("plugins") {
		if p, found := plugins[k]; found {
			log.Println("Initializing plugin", k)
			p.Init(intake.MsgBroker)
		}
	}

	// Start the HTTP server
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", viper.GetString("port")),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      intake.Router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func initConfig() {
	viper.SetDefault("port", "8080")

	viper.SetConfigName("threadle.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}
