package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/workflow/actions"
	"github.com/cryptkeeperhq/cryptkeeper/internal/workflow/decisions"
	"github.com/tidwall/gjson"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	StartNode   = "startNode"
	DefaultNode = "defaultNode"
)

type Temporal struct {
}

func (t *Temporal) GetWorkflowHistory(workflow models.Workflow) error {
	c, err := client.Dial(client.Options{})
	iter := c.GetWorkflowHistory(context.Background(), workflow.ID, "", false, 0)

	// events := []*shared.HistoryEvent{}
	for iter.HasNext() {
		event, err := iter.Next()
		if err != nil {
			return err
		}
		fmt.Println(event)
		fmt.Println("==========")
		// events = append(events, event)
	}

	return err
}

func (t *Temporal) TestPlugin(ctx context.Context, data map[string]interface{}) (string, error) {
	log.Printf("Executing TestPlugin \n\n")
	return "executed TestPlugin", nil
}

func (t *Temporal) ExecuteNode(node models.Node, context map[string]interface{}) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	switch node.Type {
	case StartNode:
		return context, nil
	case DefaultNode:
		switch node.Data.Type {
		case "decision":
			result, _ := decisions.ExecuteDecisionNode(node, context)
			if !result {
				return response, errors.New("stopping workflow execution due to false result in DecisionNode")
			}
		case "action":
			result, err := actions.ExecuteActionNode(node, context)
			return result, err
		default:
			return response, fmt.Errorf("nothing to execute for defaultNode with Type: %s, Event: %s, ID: %s", node.Data.Type, node.Data.Event, node.ID)
		}
	default:
		return response, fmt.Errorf("unsupported node type: %s", node.Data.Type)
	}

	return response, nil
}

func (t *Temporal) Execute(workflowJson models.Workflow, input map[string]interface{}) (string, error) {

	taskQueueName := "test"

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		return "", err
	}

	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        workflowJson.ID,
		TaskQueue: taskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, t.RunWorkflow, workflowJson, input)
	if err != nil {
		return "", err
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())
	// http://localhost:8080/namespaces/default/workflows/3ebc76f6-bb60-4c53-862a-c41ede651e3a/e329b757-02fe-490c-9992-3eefc6ffb9d0/history

	// var result string
	// err = we.Get(context.Background(), &result)

	// if err != nil {
	// 	return "", err
	// }

	return we.GetRunID(), nil
}

func (t *Temporal) RunWorkflow(ctx workflow.Context, workflowJson models.Workflow, input map[string]interface{}) (string, error) {

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        1, // unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	// Build a map of nodes for quick access
	nodeMap := make(map[string]models.Node)
	for _, node := range workflowJson.Details.Nodes {
		nodeMap[node.ID] = node
	}

	var workflowOutput = make(map[string]map[string]interface{})
	var mu sync.Mutex // Mutex to protect concurrent access to workflowOutput

	startNodeID := "1" // Replace with your actual start node ID
	if err := t.traverseGraph(ctx, startNodeID, workflowJson, workflowOutput, input, &mu); err != nil {
		log.Printf("Workflow execution failed: %s\n", err.Error())
	}

	// Return the output of the workflow, e.g., workflowOutput["node-4"]
	// You can customize this based on your specific use case
	return "", nil
}

func (t *Temporal) traverseGraph(ctx workflow.Context, nodeID string, workflowJson models.Workflow, workflowOutput map[string]map[string]interface{}, input map[string]interface{}, mu *sync.Mutex) error {
	node := t.getNodeByID(nodeID, workflowJson.Details.Nodes)

	t.ReplaceInputValues(node.Data.Inputs, workflowOutput)

	nodeResponse := make(map[string]interface{})
	executeNodeErr := workflow.ExecuteActivity(ctx, t.ExecuteNode, node, input).Get(ctx, &nodeResponse)
	outputValues := t.ExtractOutputValues(nodeResponse, node.Data.Outputs)

	// Protect concurrent access to workflowOutput
	mu.Lock()
	workflowOutput[node.ID] = outputValues
	mu.Unlock()

	if executeNodeErr != nil {
		log.Printf("Found error when executing node [%s]\n", executeNodeErr.Error())
		return executeNodeErr
	}

	for _, edge := range workflowJson.Details.Edges {
		if edge.Source == nodeID {
			if err := t.traverseGraph(ctx, edge.Target, workflowJson, workflowOutput, input, mu); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Temporal) getNodeByID(nodeID string, nodes []models.Node) models.Node {
	for _, node := range nodes {
		if node.ID == nodeID {
			return node
		}
	}
	return models.Node{}
}

func (t *Temporal) ReplaceInputValues(input []models.NodeInput, workflowOutput map[string]map[string]interface{}) {
	// Define a regular expression to match variable references like {{node.email}}
	regex := regexp.MustCompile(`{{(.*?)}}`)

	// Iterate through the input map
	for key, value := range input {
		if stringValue, isString := value.Value.(string); isString {
			// Find all matches of variable references in the input string
			matches := regex.FindAllString(stringValue, -1)

			// Iterate through the matches
			for _, match := range matches {
				// Extract the reference, e.g., "{{node.email}}"
				reference := match[2 : len(match)-2]

				// Split the reference into parts using "."
				parts := strings.Split(reference, ".")

				if len(parts) == 2 {
					// Get the node ID and output variable name from the reference
					nodeID := parts[0]
					outputVariable := parts[1]

					// Check if the nodeID exists in the workflowOutput map
					if outputValues, ok := workflowOutput[nodeID]; ok {
						// Check if the output variable exists in the outputValues map
						if outputValue, ok := outputValues[outputVariable]; ok {
							// Replace the variable reference with the corresponding output value
							stringValue = strings.Replace(stringValue, match, fmt.Sprintf("%v", outputValue), -1)
						}
					}
				}
			}

			// Update the input value with the modified string
			input[key].Value = stringValue
		}
	}
}

func (t *Temporal) ExtractOutputValues(nodeResponseMap map[string]interface{}, outputConfig []models.NodeOutput) map[string]interface{} {
	outputValues := make(map[string]interface{})

	nodeResponse, _ := json.Marshal(nodeResponseMap)
	for _, config := range outputConfig {
		jsonPath := config.JsonPath
		if jsonPath == "" {
			jsonPath = config.ID
		}
		variableName := config.ID

		// Use gjson to extract the value based on the JSONPath from nodeResponse
		value := gjson.Get(string(nodeResponse), jsonPath).Value()

		// Store the extracted value in the outputValues map
		outputValues[variableName] = value
	}

	return outputValues
}

// Define the topological sort function
func (t *Temporal) topologicalSort(graph map[string][]string) ([]string, error) {
	sortedNodes := []string{}
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	var visit func(node string) error
	visit = func(node string) error {
		if visited[node] {
			return nil
		}

		stack[node] = true
		for _, neighbor := range graph[node] {
			if stack[neighbor] {
				return errors.New("cycle detected in workflow graph")
			}
			if !visited[neighbor] {
				if err := visit(neighbor); err != nil {
					return err
				}
			}
		}
		stack[node] = false
		visited[node] = true
		sortedNodes = append(sortedNodes, node)
		return nil
	}

	for node := range graph {
		if !visited[node] {
			if err := visit(node); err != nil {
				return nil, err
			}
		}
	}

	// Reverse the sorted nodes to get the correct order
	for i, j := 0, len(sortedNodes)-1; i < j; i, j = i+1, j-1 {
		sortedNodes[i], sortedNodes[j] = sortedNodes[j], sortedNodes[i]
	}

	return sortedNodes, nil
}
