package main

import "time"

type Event struct {
	EventID    string            `json:"event_id"`
	CampaignID string            `json:"campaign_id"`
	ContactID  string            `json:"contact_id"`
	Type       string            `json:"type"`
	Timestamp  string            `json:"timestamp"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type EventError struct {
	Index   int    `json:"index"`
	EventID string `json:"event_id,omitempty"`
	Error   string `json:"error"`
}

type Stats struct {
	CampaignID string `json:"campaign_id"`
	Sent       int    `json:"sent"`
	Delivered  int    `json:"delivered"`
	Opened     int    `json:"opened"`
	Clicked    int    `json:"clicked"`
	Total      int    `json:"total"`
}

func validateEvent(e Event) error {
	if e.EventID == "" {
		return errString("event_id is required")
	}
	if e.CampaignID == "" {
		return errString("campaign_id is required")
	}
	if e.ContactID == "" {
		return errString("contact_id is required")
	}

	switch e.Type {
	case "sent", "delivered", "opened", "clicked":
	default:
		return errString("type must be sent, delivered, opened, or clicked")
	}

	if _, err := time.Parse(time.RFC3339, e.Timestamp); err != nil {
		return errString("timestamp must be a valid RFC3339 timestamp")
	}

	return nil
}

type errString string

func (e errString) Error() string {
	return string(e)
}
