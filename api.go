package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

const createEndpoint = `https://api.twitch.tv/helix/eventsub/subscriptions`

var (
	client = http.Client{
		Timeout: time.Second * 5,
	}
)

type Credentials struct {
	ClientID    string
	AccessToken string
}

type CreateParams struct {
	SubscriptionType     string
	SubscriptionVersions string
	Condition            json.RawMessage
}

type createBody struct {
	SubscriptionType    string          `json:"type"`
	SubscriptionVersion string          `json:"version"`
	Condition           json.RawMessage `json:"condition"`
	Transport           Transport       `json:"transport"`
}

func CreateWebsocketSubscription(sessionID string, params CreateParams, creds Credentials) error {
	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(createBody{
		SubscriptionType:    params.SubscriptionType,
		SubscriptionVersion: params.SubscriptionVersions,
		Condition:           params.Condition,
		Transport: Transport{
			Method:    "websocket",
			SessionID: sessionID,
		},
	})
	if err != nil {
		return err
	}

	log.Info(buf.String())

	req, err := http.NewRequest("POST", createEndpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", creds.AccessToken))
	req.Header.Set("Client-Id", creds.ClientID)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("bad status code: %d %s", res.StatusCode, string(body))
	}

	return nil
}
