package models

type Event struct {
	EventType string `json:"eventType"`
	User      string `json:"user"`
}
