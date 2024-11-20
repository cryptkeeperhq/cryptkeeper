package actions

import "github.com/cryptkeeperhq/cryptkeeper/internal/models"

type EmailActionPlugin struct {
	// You can add any configuration or dependencies needed for the plugin here
}

func (e *EmailActionPlugin) Execute(inputs []models.NodeInput) (map[string]interface{}, error) {
	var to, subject, message string
	for _, input := range inputs {
		if input.Value == nil {
			continue
		}
		switch input.ID {
		case "to":
			to = input.Value.(string)
		case "subject":
			subject = input.Value.(string)
		case "message":
			message = input.Value.(string)
		}
	}

	// Return any relevant results or error
	result := map[string]interface{}{
		"status":  "Email sent",
		"to":      to,
		"subject": subject,
		"message": message,
	}
	return result, nil
}
