package ingestor

import (
  "net/http"
)

type EventProcessor struct {
  db Database
}

func (h *EventProcessor) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hook", collectorHandler())
	return mux
}
