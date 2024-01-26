package main

import (
	"encoding/json"
	"time"
)

type Message struct {
	Metadata MessageMetadata `json:"metadata"`
	Payload  json.RawMessage `json:"payload"`
}

type Transport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

type MessageMetadata struct {
	MessageID        string    `json:"message_id"`
	MessageType      string    `json:"message_type"`
	MessageTimestamp time.Time `json:"message_timestamp"`

	SubscriptionType    string `json:"subscription_type,omitempty"`
	SubscriptionVersion string `json:"subscription_version,omitempty"`
}

type Session struct {
	ID                      string    `json:"id"`
	Status                  string    `json:"status"`
	ConnectedAt             time.Time `json:"connected_at"`
	KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
	ReconnectURL            *string   `json:"reconnect_url,omitempty"`
}

type SessionWelcomePayload struct {
	Session Session `json:"session"`
}

type Subscription struct {
	ID        string          `json:"id"`
	Status    string          `json:"status"`
	Type      string          `json:"type"`
	Version   string          `json:"version"`
	Cost      int             `json:"cost"`
	Condition json.RawMessage `json:"condition"`
	Transport Transport       `json:"transport"`
	CreatedAt time.Time       `json:"created_at"`
}

type NotificationPayload struct {
	Subscription Subscription    `json:"subscription"`
	Event        json.RawMessage `json:"event"`
}

type ReconnectPayload struct {
	Session Session `json:"session"`
}

type RevocationPayload struct {
	Subscription Subscription `json:"subscription"`
}
