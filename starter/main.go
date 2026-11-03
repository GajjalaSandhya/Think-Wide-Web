package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"sync"
)

var (
	storeMu sync.RWMutex
	events  = make(map[string]Event)
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /events", handlePostEvents)
	mux.HandleFunc("GET /campaigns/{campaignID}/stats", handleGetStats)

	log.Println("Campaign Events API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handlePostEvents(w http.ResponseWriter, r *http.Request) {
	const maxBodySize = 5 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "unable to read request body",
		})
		return
	}

	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "request body cannot be empty",
		})
		return
	}

	var rawEvents []json.RawMessage
	if err := json.Unmarshal(body, &rawEvents); err != nil || rawEvents == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "request body must be a JSON array of events",
		})
		return
	}

	accepted := 0
	duplicates := 0
	rejected := make([]EventError, 0)

	for i, raw := range rawEvents {
		var event Event

		if err := json.Unmarshal(raw, &event); err != nil {
			rejected = append(rejected, EventError{
				Index: i,
				Error: "invalid event JSON: " + err.Error(),
			})
			continue
		}

		if err := validateEvent(event); err != nil {
			rejected = append(rejected, EventError{
				Index:   i,
				EventID: event.EventID,
				Error:   err.Error(),
			})
			continue
		}

		storeMu.Lock()
		existing, found := events[event.EventID]

		if found {
			if !reflect.DeepEqual(existing, event) {
				storeMu.Unlock()
				rejected = append(rejected, EventError{
					Index:   i,
					EventID: event.EventID,
					Error:   "event_id already exists with different data",
				})
				continue
			}

			duplicates++
			storeMu.Unlock()
			continue
		}

		events[event.EventID] = event
		accepted++
		storeMu.Unlock()
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"accepted":   accepted,
		"duplicates": duplicates,
		"rejected":   len(rejected),
		"errors":     rejected,
	})
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {
	campaignID := r.PathValue("campaignID")

	stats := Stats{
		CampaignID: campaignID,
	}

	storeMu.RLock()
	defer storeMu.RUnlock()

	for _, event := range events {
		if event.CampaignID != campaignID {
			continue
		}

		stats.Total++

		switch event.Type {
		case "sent":
			stats.Sent++
		case "delivered":
			stats.Delivered++
		case "opened":
			stats.Opened++
		case "clicked":
			stats.Clicked++
		}
	}

	writeJSON(w, http.StatusOK, stats)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}
