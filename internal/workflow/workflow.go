package workflow

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"regexp"
// 	"strings"
// 	"sync"

// 	"github.com/google/uuid"
// 	"github.com/tidwall/gjson"
// 	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
// 	"github.com/cryptkeeperhq/cryptkeeper/internal/workflow/actions"
// 	"github.com/cryptkeeperhq/cryptkeeper/internal/workflow/decisions"
// )

// // func TestPlugin(ctx context.Context, data map[string]interface{}) (string, error) {
// // 	log.Printf("Executing TestPlugin \n\n")
// // 	return "executed TestPlugin", nil
// // }

// // func ExecuteNode(node models.Node, context map[string]interface{}) (map[string]interface{}, error) {
// // 	response := make(map[string]interface{})

// // 	switch node.Type {
// // 	case "startNode":
// // 		return context, nil
// // 	case "defaultNode":
// // 		switch node.Data.Type {
// // 		case "decision":
// // 			result, _ := decisions.ExecuteDecisionNode(node, context)
// // 			if !result {
// // 				return response, errors.New("stopping workflow execution due to false result in DecisionNode")
// // 			}
// // 		case "action":
// // 			result, err := actions.ExecuteActionNode(node, context)
// // 			return result, err
// // 		default:
// // 			return response, fmt.Errorf("nothing to execute for defaultNode with Type: %s, Event: %s, ID: %s", node.Data.Type, node.Data.Event, node.ID)
// // 		}
// // 	default:
// // 		return response, fmt.Errorf("unsupported node type: %s", node.Data.Type)
// // 	}

// // 	return response, nil
// // }

// // func Execute(request models.Workflow, taskQueueName string, input map[string]interface{}) (string, error) {

// // 	runID := uuid.New().String()
// // 	log.Printf("WorkflowID: %s RunID: %s\n", request.ID, runID)

// // 	o, e := RunLocalWorkflow(request, input)
// // 	return o, e
// // }

// // func RunLocalWorkflow(request models.Workflow, input map[string]interface{}) (string, error) {
// // 	// Build a map of nodes for quick access
// // 	nodeMap := make(map[string]models.Node)
// // 	for _, node := range request.Details.Nodes {
// // 		nodeMap[node.ID] = node
// // 	}

// // 	var workflowOutput = make(map[string]map[string]interface{})
// // 	var mu sync.Mutex // Mutex to protect concurrent access to workflowOutput

// // 	startNodeID := "1" // Replace with your actual start node ID
// // 	if err := traverseGraph(context.Background(), startNodeID, request, workflowOutput, input, &mu); err != nil {
// // 		log.Printf("Workflow execution failed: %s\n", err.Error())
// // 	}

// // 	// Return the output of the workflow, e.g., workflowOutput["node-4"]
// // 	// You can customize this based on your specific use case
// // 	return "", nil
// // }

// // func traverseGraph(ctx context.Context, nodeID string, request models.Workflow, workflowOutput map[string]map[string]interface{}, input map[string]interface{}, mu *sync.Mutex) error {
// // 	node := getNodeByID(nodeID, request.Details.Nodes)

// // 	ReplaceInputValues(node.Data.Inputs, workflowOutput)

// // 	fmt.Println(node, input)
// // 	nodeResponse, executeNodeErr := ExecuteNode(node, input)
// // 	fmt.Println(nodeResponse, executeNodeErr)
// // 	outputValues := ExtractOutputValues(nodeResponse, node.Data.Outputs)

// // 	// Protect concurrent access to workflowOutput
// // 	mu.Lock()
// // 	workflowOutput[node.ID] = outputValues
// // 	mu.Unlock()

// // 	if executeNodeErr != nil {
// // 		log.Printf("Found error when executing node [%s]\n", executeNodeErr.Error())
// // 		return executeNodeErr
// // 	}

// // 	for _, edge := range request.Details.Edges {
// // 		if edge.Source == nodeID {
// // 			if err := traverseGraph(ctx, edge.Target, request, workflowOutput, input, mu); err != nil {
// // 				return err
// // 			}
// // 		}
// // 	}

// // 	return nil
// // }

// // func getNodeByID(nodeID string, nodes []models.Node) models.Node {
// // 	for _, node := range nodes {
// // 		if node.ID == nodeID {
// // 			return node
// // 		}
// // 	}
// // 	return models.Node{}
// // }

// // func ReplaceInputValues(input []models.NodeInput, workflowOutput map[string]map[string]interface{}) {
// // 	// Define a regular expression to match variable references like {{node.email}}
// // 	regex := regexp.MustCompile(`{{(.*?)}}`)

// // 	// Iterate through the input map
// // 	for key, value := range input {
// // 		if stringValue, isString := value.Value.(string); isString {
// // 			// Find all matches of variable references in the input string
// // 			matches := regex.FindAllString(stringValue, -1)

// // 			// Iterate through the matches
// // 			for _, match := range matches {
// // 				// Extract the reference, e.g., "{{node.email}}"
// // 				reference := match[2 : len(match)-2]

// // 				// Split the reference into parts using "."
// // 				parts := strings.Split(reference, ".")

// // 				if len(parts) == 2 {
// // 					// Get the node ID and output variable name from the reference
// // 					nodeID := parts[0]
// // 					outputVariable := parts[1]

// // 					// Check if the nodeID exists in the workflowOutput map
// // 					if outputValues, ok := workflowOutput[nodeID]; ok {
// // 						// Check if the output variable exists in the outputValues map
// // 						if outputValue, ok := outputValues[outputVariable]; ok {
// // 							// Replace the variable reference with the corresponding output value
// // 							stringValue = strings.Replace(stringValue, match, fmt.Sprintf("%v", outputValue), -1)
// // 						}
// // 					}
// // 				}
// // 			}

// // 			// Update the input value with the modified string
// // 			input[key].Value = stringValue
// // 		}
// // 	}
// // }

// // func ExtractOutputValues(nodeResponseMap map[string]interface{}, outputConfig []models.NodeOutput) map[string]interface{} {
// // 	outputValues := make(map[string]interface{})

// // 	nodeResponse, _ := json.Marshal(nodeResponseMap)
// // 	for _, config := range outputConfig {
// // 		jsonPath := config.JsonPath
// // 		if jsonPath == "" {
// // 			jsonPath = config.ID
// // 		}
// // 		variableName := config.ID

// // 		// Use gjson to extract the value based on the JSONPath from nodeResponse
// // 		value := gjson.Get(string(nodeResponse), jsonPath).Value()

// // 		// Store the extracted value in the outputValues map
// // 		outputValues[variableName] = value
// // 	}

// // 	return outputValues
// // }

// // // Define the topological sort function
// // func topologicalSort(graph map[string][]string) ([]string, error) {
// // 	sortedNodes := []string{}
// // 	visited := make(map[string]bool)
// // 	stack := make(map[string]bool)

// // 	var visit func(node string) error
// // 	visit = func(node string) error {
// // 		if visited[node] {
// // 			return nil
// // 		}

// // 		stack[node] = true
// // 		for _, neighbor := range graph[node] {
// // 			if stack[neighbor] {
// // 				return errors.New("cycle detected in workflow graph")
// // 			}
// // 			if !visited[neighbor] {
// // 				if err := visit(neighbor); err != nil {
// // 					return err
// // 				}
// // 			}
// // 		}
// // 		stack[node] = false
// // 		visited[node] = true
// // 		sortedNodes = append(sortedNodes, node)
// // 		return nil
// // 	}

// // 	for node := range graph {
// // 		if !visited[node] {
// // 			if err := visit(node); err != nil {
// // 				return nil, err
// // 			}
// // 		}
// // 	}

// // 	// Reverse the sorted nodes to get the correct order
// // 	for i, j := 0, len(sortedNodes)-1; i < j; i, j = i+1, j-1 {
// // 		sortedNodes[i], sortedNodes[j] = sortedNodes[j], sortedNodes[i]
// // 	}

// // 	return sortedNodes, nil
// // }
