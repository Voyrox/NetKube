package pod

import "netkube/adapters/api/shared"

type Row struct {
	Namespace         string `json:"namespace"`
	Name              string `json:"name"`
	Ready             string `json:"ready"`
	Status            string `json:"status"`
	Restarts          int32  `json:"restarts"`
	LastRestart       string `json:"lastRestart"`
	LastRestartReason string `json:"lastRestartReason"`
	Node              string `json:"node"`
	PodIP             string `json:"podIP"`
	Age               string `json:"age"`
}

type Stats struct {
	Running int `json:"running"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
	Other   int `json:"other"`
	Total   int `json:"total"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []Row           `json:"items"`
	Count int             `json:"count"`
	Stats Stats           `json:"stats"`
}

type ContainerRow struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Ready    bool   `json:"ready"`
	Restarts int32  `json:"restarts"`
	State    string `json:"state"`
}

type ConditionRow struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type Detail struct {
	Namespace         string            `json:"namespace"`
	Name              string            `json:"name"`
	Ready             string            `json:"ready"`
	Status            string            `json:"status"`
	Phase             string            `json:"phase"`
	Restarts          int32             `json:"restarts"`
	LastRestart       string            `json:"lastRestart"`
	LastRestartReason string            `json:"lastRestartReason"`
	Node              string            `json:"node"`
	PodIP             string            `json:"podIP"`
	HostIP            string            `json:"hostIP"`
	ServiceAccount    string            `json:"serviceAccount"`
	QOSClass          string            `json:"qosClass"`
	Age               string            `json:"age"`
	StartTime         string            `json:"startTime"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	Containers        []ContainerRow    `json:"containers"`
	Conditions        []ConditionRow    `json:"conditions"`
}

type DetailResponse struct {
	Meta shared.PageMeta `json:"meta"`
	Item Detail          `json:"item"`
}

type LogResponse struct {
	Meta      shared.PageMeta `json:"meta"`
	Container string          `json:"container"`
	Content   string          `json:"content"`
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
