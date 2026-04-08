package deployment

import "netkube/adapters/api/shared"

type Row struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Ready     string `json:"ready"`
	Status    string `json:"status"`
	Desired   int32  `json:"desired"`
	Updated   int32  `json:"updated"`
	Available int32  `json:"available"`
	Age       string `json:"age"`
}

type Stats struct {
	Healthy int `json:"healthy"`
	Warning int `json:"warning"`
	Pending int `json:"pending"`
	Total   int `json:"total"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []Row           `json:"items"`
	Count int             `json:"count"`
	Error string          `json:"error,omitempty"`
	Stats Stats           `json:"stats"`
}

type ConditionRow struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type ReplicaSetRow struct {
	Name    string `json:"name"`
	Ready   string `json:"ready"`
	Desired int32  `json:"desired"`
	Age     string `json:"age"`
}

type PodRow struct {
	Name   string `json:"name"`
	Ready  string `json:"ready"`
	Status string `json:"status"`
	Node   string `json:"node"`
	Age    string `json:"age"`
}

type Detail struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Ready       string            `json:"ready"`
	Desired     int32             `json:"desired"`
	Updated     int32             `json:"updated"`
	Available   int32             `json:"available"`
	Unavailable int32             `json:"unavailable"`
	Age         string            `json:"age"`
	Strategy    string            `json:"strategy"`
	Selector    string            `json:"selector"`
	Conditions  []ConditionRow    `json:"conditions"`
	ReplicaSets []ReplicaSetRow   `json:"replicaSets"`
	Pods        []PodRow          `json:"pods"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type DetailResponse struct {
	Meta shared.PageMeta `json:"meta"`
	Item Detail          `json:"item"`
}

type EventRow struct {
	Type    string `json:"type"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Age     string `json:"age"`
}

type EventsResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []EventRow      `json:"items"`
}

type YAMLResponse struct {
	Meta    shared.PageMeta `json:"meta"`
	Content string          `json:"content"`
}
