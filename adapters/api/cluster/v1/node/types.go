package node

import "netkube/adapters/api/shared"

type Summary struct {
	Total    int    `json:"total"`
	Ready    int    `json:"ready"`
	Pending  int    `json:"pending"`
	NotReady int    `json:"notReady"`
	Status   string `json:"status"`
}

type ConditionRow struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type EventRow struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Age     string `json:"age"`
	Type    string `json:"type"`
}

type ListItem struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Role       string `json:"role"`
	Version    string `json:"version"`
	InternalIP string `json:"internalIP"`
	Age        string `json:"age"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}

type Detail struct {
	Name              string            `json:"name"`
	Status            string            `json:"status"`
	Role              string            `json:"role"`
	KubeletVersion    string            `json:"kubeletVersion"`
	ContainerRuntime  string            `json:"containerRuntime"`
	OSKernel          string            `json:"osKernel"`
	Architecture      string            `json:"architecture"`
	InternalIP        string            `json:"internalIP"`
	PodCIDR           string            `json:"podCIDR"`
	AllocatableCPU    string            `json:"allocatableCPU"`
	AllocatableMemory string            `json:"allocatableMemory"`
	AllocatablePods   string            `json:"allocatablePods"`
	AllocatableStore  string            `json:"allocatableStore"`
	CPUPressure       string            `json:"cpuPressure"`
	MemoryPressure    string            `json:"memoryPressure"`
	DiskPressure      string            `json:"diskPressure"`
	Labels            map[string]string `json:"labels"`
	Conditions        []ConditionRow    `json:"conditions"`
	Timeline          []EventRow        `json:"timeline"`
}

type DetailResponse struct {
	Meta shared.PageMeta `json:"meta"`
	Item Detail          `json:"item"`
}

type YAMLResponse struct {
	Meta    shared.PageMeta `json:"meta"`
	Content string          `json:"content"`
}
