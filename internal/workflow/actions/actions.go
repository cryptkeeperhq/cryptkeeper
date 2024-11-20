package actions

import (
	"fmt"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

type ActionPlugin interface {
	Execute(inputs []models.NodeInput) (map[string]interface{}, error)
}

func CreateActionPlugin(event string) (ActionPlugin, error) {
	// Implement logic to create the appropriate plugin based on the event
	switch event {
	case "email":
		return &EmailActionPlugin{}, nil
	case "sms":
		return &SMSActionPlugin{}, nil
	case "push-notification":
		return &PushNotificationActionPlugin{}, nil
	case "api":
		return &APIActionPlugin{}, nil
	case "slack":
		return &SlackActionPlugin{}, nil
	default:
		return nil, fmt.Errorf("unsupported action event: %s", event)
	}
}

func ExecuteActionNode(node models.Node, context map[string]interface{}) (map[string]interface{}, error) {
	if node.Data.Type != "action" {
		return nil, fmt.Errorf("unsupported node type: %s", node.Data.Type)
	}

	// Create the corresponding action plugin
	plugin, err := CreateActionPlugin(node.Data.Event)
	if err != nil {
		return nil, err
	}

	// Execute the action using the plugin
	result, err := plugin.Execute(node.Data.Inputs)
	if err != nil {
		return nil, err
	}

	return result, nil
}
