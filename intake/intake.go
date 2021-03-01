package intake

import (
	"context"
	fmt "fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// Endpoint name constants
const (
	IntakeEndpointV1 = "/intake/"

	SeriesEndpointV1       = "/api/v1/series"
	CheckRunsEndpointV1    = "/api/v1/check_run"
	SketchSeriesEndpointV1 = "/api/v1/sketches"
	ValidateEndpointV1     = "/api/v1/validate"
	ProcessesEndpointV1    = "/api/v1/collector"
	ContainerEndpointV1    = "/api/v1/container"
	OrchestratorEndpointV1 = "/api/v1/orchestrator"

	SeriesEndpointV2        = "/api/v2/series"
	EventsEndpointV2        = "/api/v2/events"
	ServiceChecksEndpointV2 = "/api/v2/service_checks"
	SketchSeriesEndpointV2  = "/api/beta/sketches"
	HostMetadataEndpointV2  = "/api/v2/host_metadata"
	MetadataEndpointV2      = "/api/v2/metadata"

	v1PathPrefix = "/api/v1"
)

var (
	// MsgBroker is used to read messages from the intake
	MsgBroker *PubSub
	// router is the API router
	router *mux.Router
)

// Init the message broker and the API router
func init() {
	MsgBroker = NewPubsub()

	router = mux.NewRouter()
	// always validate the api key
	router.HandleFunc(ValidateEndpointV1, func(rw http.ResponseWriter, r *http.Request) {})
	// API v1, multiple endpoints
	router.PathPrefix(v1PathPrefix).HandlerFunc(v1Handler)
	// intake, single endpoint
	router.HandleFunc(IntakeEndpointV1, intakeHandler)
	// catch-all route, for debug and unsupported endpoints
	router.PathPrefix("/").HandlerFunc(defaultHandler)
}

func defaultHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
}

// This handler serves the /intake/ endpoint
func intakeHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := readRequestBody(r)
	if err != nil {
		log.Println("intakeHandler: error reading request body:", err)
		http.Error(rw, "", http.StatusBadRequest)
		return
	}
	MsgBroker.Publish(r.URL.Path, body)
}

// This handler serves the api/v1/* endpoints, deflates the request body
// and sends the paylaod to the message broker using the URL path as the
// topic name.
func v1Handler(rw http.ResponseWriter, r *http.Request) {
	body, err := readRequestBody(r)
	if err != nil {
		log.Println("v1Handler: error reading request body:", err)
		http.Error(rw, "", http.StatusBadRequest)
		return
	}
	MsgBroker.Publish(r.URL.Path, body)
}

// GetV1Endpoints returns a slice containing all the v1 endpoints
func GetV1Endpoints() []string {
	return []string{
		SeriesEndpointV1,
		CheckRunsEndpointV1,
		SketchSeriesEndpointV1,
		ValidateEndpointV1,
		ProcessesEndpointV1,
		ContainerEndpointV1,
		OrchestratorEndpointV1,
	}
}

// Serve starts the HTTP server and blocks
func Serve() {
	// Start the HTTP server
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", viper.GetString("port")),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
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
}
