package models

type Event struct {
	EventType string `json:"eventType"`
	UserID    string `json:"userID"`
	Age       int    `json:"age"`
	NoFiles   int    `json:"noFiles"`
}
