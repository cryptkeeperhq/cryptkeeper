package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

type SlackActionPlugin struct {
	WebhookURL string
}

func (s *SlackActionPlugin) Execute(inputs []models.NodeInput) (map[string]interface{}, error) {
	// Initialize variables for Slack message parameters
	var channel, messageType, messageText string

	// Extract values from inputs
	for _, input := range inputs {
		if input.Value == nil {
			continue
		}
		switch input.ID {
		case "webhookUrl":
			s.WebhookURL = input.Value.(string)
		case "channel":
			channel = input.Value.(string)
		case "messageType":
			messageType = input.Value.(string)
		case "message":
			messageText = input.Value.(string)
		}
	}

	// Determine the color based on messageType
	color := "#36a64f" // Default color (success)
	switch messageType {
	case "warning":
		color = "#ffcc00" // Yellow for warning
	case "danger":
		color = "#ff0000" // Red for danger
	}

	// Construct the Slack message payload
	slackMessage := map[string]interface{}{
		"channel":    channel,
		"username":   "Workflow Bot",
		"icon_emoji": ":robot_face:",
		"attachments": []map[string]interface{}{
			{
				"text":  messageText,
				"color": color,
			},
		},
	}

	// Send the message to Slack
	err := s.sendSlackMessage(slackMessage)
	if err != nil {
		return nil, err
	}

	// Return a success message
	result := map[string]interface{}{
		"status":  "Message sent to Slack",
		"channel": channel,
	}
	return result, nil
}

func (s *SlackActionPlugin) sendSlackMessage(message map[string]interface{}) error {
	// Marshal the message payload to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling Slack message payload: %v", err)
	}

	// Create an HTTP POST request to the Slack webhook URL
	req, err := http.NewRequest("POST", s.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set a timeout for the HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request to Slack: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack webhook returned non-OK status code: %d", resp.StatusCode)
	}

	return nil
}
