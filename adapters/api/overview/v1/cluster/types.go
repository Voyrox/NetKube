package cluster

import "netkube/adapters/api/shared"

type Metric struct {
	Total       int    `json:"total"`
	Primary     int    `json:"primary"`
	Warning     int    `json:"warning"`
	Danger      int    `json:"danger"`
	Other       int    `json:"other,omitempty"`
	HealthTotal int    `json:"healthTotal,omitempty"`
	Status      string `json:"status"`
}
type ResourceUsageMetric struct {
	Percent float64 `json:"percent"`
	Used    string  `json:"used"`
	Total   string  `json:"total"`
}
type ResourceUsageSection struct {
	CPU    ResourceUsageMetric `json:"cpu"`
	Memory ResourceUsageMetric `json:"memory"`
	Pods   ResourceUsageMetric `json:"pods,omitempty"`
}
type ResourceUsage struct {
	UsageCapacity    ResourceUsageSection `json:"usageCapacity"`
	RequestsAllocate ResourceUsageSection `json:"requestsAllocate"`
}
type Response struct {
	Meta              shared.PageMeta       `json:"meta"`
	Nodes             Metric                `json:"nodes"`
	PersistentVolumes Metric                `json:"persistentVolumes"`
	CustomResources   Metric                `json:"customResources"`
	ResourceUsage     ResourceUsage         `json:"resourceUsage"`
	Warnings          []shared.WarningEvent `json:"warnings"`
}
