package intake

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	// Router is the API router
	Router *mux.Router
)

// Init the message broker and the API router
func Init() {
	MsgBroker = NewPubsub()

	Router = mux.NewRouter()
	// always validate the api key
	Router.HandleFunc(ValidateEndpointV1, func(rw http.ResponseWriter, r *http.Request) {})
	// API v1, multiple endpoints
	Router.PathPrefix(v1PathPrefix).HandlerFunc(v1Handler)
	// intake, single endpoint
	Router.HandleFunc(IntakeEndpointV1, intakeHandler)
	// catch-all route, for debug and unsupported endpoints
	Router.PathPrefix("/").HandlerFunc(defaultHandler)
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
