package elasticsearch

import (
	"log"

	"github.com/masci/threadle/intake"
)

func readV1Metrics(broker *intake.PubSub, exclude filters) {
	go func() {
		var metrics []intake.V1Metric
		var err error
		for msg := range broker.Subscribe(intake.SeriesEndpointV1) {
			metrics, err = intake.GetV1Metrics([]byte(msg))
			if err != nil {
				log.Println("error processing metrics: ", err)
				continue
			}

			for _, m := range filterV1Metrics(metrics, exclude) {
				log.Println(m.Metric)
			}
		}
	}()
}

func filterV1Metrics(metrics []intake.V1Metric, exclude filters) []intake.V1Metric {
	out := []intake.V1Metric{}

	for _, m := range metrics {
		// drop excluded metrics
		for _, reg := range exclude {
			if reg.MatchString(m.Metric) {
				continue
			}
		}

		out = append(out, m)
	}
	return out
}
