package events

import "time"

const (
	SecretCreated = "secret.created"
	SecretUpdated = "secret.updated"
	SecretDeleted = "secret.deleted"
	SecretRotated = "secret.rotated"
	PathCreated   = "path.created"
	PathDeleted   = "path.deleted"
)

type User struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

type Event struct {
	EventType string                 `json:"eventType"`
	Timestamp time.Time              `json:"timestamp"`
	User      User                   `json:"user"`
	Details   map[string]interface{} `json:"details"`
}

type EventField struct {
	EventType string   `json:"eventType"`
	Fields    []string `json:"fields"`
}

type EventListResponse struct {
	Events []EventField `json:"events"`
}

func GetEventList() EventListResponse {
	return EventListResponse{
		Events: []EventField{
			{
				EventType: "secret.created",
				Fields:    []string{"secretId", "secretName", "timestamp", "user"},
			},
			{
				EventType: "secret.updated",
				Fields:    []string{"secretId", "secretName", "timestamp", "user"},
			},
			{
				EventType: "secret.deleted",
				Fields:    []string{"secretId", "timestamp", "user"},
			},
			{
				EventType: "secret.rotated",
				Fields:    []string{"secretId", "timestamp", "user"},
			},
			{
				EventType: "path.created",
				Fields:    []string{"path", "timestamp", "user"},
			},
			{
				EventType: "path.deleted",
				Fields:    []string{"path", "timestamp", "user"},
			},
		},
	}
}

func NewEvent(eventType string, user User, details map[string]interface{}) Event {
	return Event{
		EventType: eventType,
		Timestamp: time.Now(),
		User:      user,
		Details:   details,
	}
}

func PublishEvent() {

}
