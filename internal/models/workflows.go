package models

import (
	"time"
)

type Workflow struct {
	ID        string       `json:"id" pg:"id,pk"`
	Name      string       `json:"name" pg:"name"`
	Details   WorkflowJSON `json:"details"`
	CreatedBy string       `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkflowJSON struct {
	Event   map[string]interface{} `json:"event"`
	Trigger string                 `json:"trigger"`
	Nodes   []Node                 `json:"nodes"`
	Edges   []Edge                 `json:"edges"`
}

type Node struct {
	Deletable bool   `json:"deletable"`
	ID        string `json:"id"`
	Type      string `json:"type"`
	Position  struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"`
	Data             NodeData `json:"data,omitempty"`
	Width            int      `json:"width"`
	Height           int      `json:"height"`
	Selected         bool     `json:"selected,omitempty"`
	PositionAbsolute struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"positionAbsolute,omitempty"`
	Dragging bool `json:"dragging,omitempty"`
}

type NodeData struct {
	Label            string                   `json:"label"`
	Icon             string                   `json:"icon"`
	Type             string                   `json:"type"`
	Event            string                   `json:"event"`
	Description      string                   `json:"description"`
	BorderColor      string                   `json:"borderColor"`
	Inputs           []NodeInput              `json:"inputs"`
	Outputs          []NodeOutput             `json:"outputs"`
	AttributeConfigs []map[string]interface{} `json:"attributeConfigs"`
}

type NodeInput struct {
	ID        string      `json:"id"`
	Label     string      `json:"label"`
	InputType string      `json:"inputType"`
	Values    []string    `json:"values,omitempty"`
	Value     interface{} `json:"value"`
}

type NodeOutput struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	JsonPath string `json:"jsonPath"`
	DataType string `json:"dataType"`
}

type Edge struct {
	Type      string `json:"type"`
	MarkerEnd struct {
		Type string `json:"type"`
	} `json:"markerEnd"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	Target       string `json:"target"`
	TargetHandle string `json:"targetHandle"`
	ID           string `json:"id"`
}
