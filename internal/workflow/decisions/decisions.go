package decisions

import (
	"fmt"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
)

func ExecuteDecisionNode(node models.Node, context map[string]interface{}) (bool, error) {
	// Implement the logic for "decision" type nodes
	fmt.Printf("Executing DecisionNode with Type: %s, Event: %s, ID: %s\n", node.Data.Type, node.Data.Event, node.ID)

	// Initialize a result variable
	result := true

	// Now you can work with the attributeConfigs in your logic
	for _, attrConfig := range node.Data.AttributeConfigs {
		attribute := attrConfig["attribute"].(string)
		operator := attrConfig["operator"].(string)
		value := attrConfig["value"].(string)

		// Access the attribute value from the context based on the attribute name
		attrValue, ok := context[attribute].(string)
		if !ok {
			return false, fmt.Errorf("attribute '%s' not found in context", attribute)
		}

		// Perform comparisons or other logic using attribute, operator, and value
		switch operator {
		case "==":
			if attrValue != value {
				result = false
			}
		case "!=":
			if attrValue == value {
				result = false
			}
		// Add more comparison cases as needed
		default:
			return false, fmt.Errorf("unsupported operator: %s", operator)
		}
	}

	return result, nil
}
