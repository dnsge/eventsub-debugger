package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"
	basicws "github.com/dnsge/go-basic-websocket"
)

var (
	server    = flag.String("server", "wss://eventsub.wss.twitch.tv/ws", "Twitch EventSub WebSocket endpoint")
	clientID  = flag.String("client-id", "", "Twitch Application ClientID")
	authToken = flag.String("token", "", "Twitch OAuth Access Token")

	subscriptionType    = flag.String("sub-type", "", "EventSub Subscription Type")
	subscriptionVersion = flag.String("sub-version", "", "EventSub Subscription Version")
	condition           = flag.String("condition", "", "JSON Condition to pass")
)

func init() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
}

func main() {
	if *clientID == "" {
		log.Fatal("Client ID must be set")
	}
	if *authToken == "" {
		log.Fatal("Auth token must be set")
	}
	if *subscriptionType == "" {
		log.Fatal("Subscription type must be set")
	}
	if *subscriptionVersion == "" {
		log.Fatal("Subscription version must be set")
	}
	if *condition == "" {
		log.Fatal("Condition version must be set")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ws := basicws.NewBasicWebsocket(*server, http.Header{})
	s := &session{}

	ws.AutoReconnect = false

	ws.OnConnect = func() {
		log.Info("Connected")
	}

	ws.OnMessage = func(b []byte) error {
		log.Debug("[RECV] " + string(b))

		var data Message
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
		s.processMessage(&data)
		return nil
	}

	ws.OnError = func(err error) {
		log.Error(err.Error())
	}

	// Connect and wait for interrupt
	err := ws.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-interrupt
}

type session struct {
	sessionID string
}

func (s *session) processMessage(msg *Message) {
	switch msg.Metadata.MessageType {
	case "session_welcome":
		var payload SessionWelcomePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Error("failed to unmarshal session_welcome", "err", err)
			return
		}
		s.processWelcome(msg, &payload)
	case "session_keepalive":
		s.processKeepalive(msg)
	case "notification":
		var payload NotificationPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Error("failed to unmarshal notification", "err", err)
			return
		}
		s.processNotification(msg, &payload)
	case "session_reconnect":
		var payload ReconnectPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Error("failed to unmarshal session_reconnect", "err", err)
			return
		}
		s.processReconnect(msg, &payload)
	case "revocation":
		var payload RevocationPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Error("failed to unmarshal revocation", "err", err)
			return
		}
		s.processRevocation(msg, &payload)
	}
}

func (s *session) processWelcome(msg *Message, payload *SessionWelcomePayload) {
	log.Info("session_welcome", "message_id", msg.Metadata.MessageID, "session_id", payload.Session.ID)
	s.sessionID = payload.Session.ID

	err := CreateWebsocketSubscription(s.sessionID, CreateParams{
		SubscriptionType:     *subscriptionType,
		SubscriptionVersions: *subscriptionVersion,
		Condition:            []byte(*condition),
	}, Credentials{
		ClientID:    *clientID,
		AccessToken: *authToken,
	})
	if err != nil {
		log.Error("failed to create websocket subscription", "err", err)
		os.Exit(1)
	}
}

func (s *session) processKeepalive(msg *Message) {
	log.Info("session_keepalive", "message_id", msg.Metadata.MessageID)
}

func (s *session) processNotification(msg *Message, payload *NotificationPayload) {
	subscription := fmt.Sprintf("%s.%s", msg.Metadata.SubscriptionType, msg.Metadata.SubscriptionVersion)
	log.Info("notification", "message_id", msg.Metadata.MessageID, "subscription", subscription)

	var eventData map[string]any
	if err := json.Unmarshal(payload.Event, &eventData); err != nil {
		log.Error("failed to unmarshal event", "err", err)
		return
	}

	pretty, _ := json.MarshalIndent(eventData, "", "  ")
	log.Info(string(pretty))
}

func (s *session) processReconnect(msg *Message, payload *ReconnectPayload) {
	log.Info("session_reconnect", "message_id", msg.Metadata.MessageID)
}

func (s *session) processRevocation(msg *Message, payload *RevocationPayload) {
	log.Info("revocation", "message_id", msg.Metadata.MessageID)
	pretty, _ := json.MarshalIndent(payload, "", "  ")
	log.Info(pretty)
}
