package event

import "netkube/adapters/api/shared"

type TimelineRow struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Age     string `json:"age"`
	Type    string `json:"type"`
}

type Detail struct {
	Title          string            `json:"title"`
	Type           string            `json:"type"`
	Namespace      string            `json:"namespace"`
	Reason         string            `json:"reason"`
	InvolvedObject string            `json:"involvedObject"`
	Kind           string            `json:"kind"`
	Name           string            `json:"name"`
	Source         string            `json:"source"`
	FirstSeen      string            `json:"firstSeen"`
	LastSeen       string            `json:"lastSeen"`
	Count          int32             `json:"count"`
	Node           string            `json:"node"`
	Message        string            `json:"message"`
	Timeline       []TimelineRow     `json:"timeline"`
	Annotations    map[string]string `json:"annotations"`
}

type DetailResponse struct {
	Meta shared.PageMeta `json:"meta"`
	Item Detail          `json:"item"`
}
