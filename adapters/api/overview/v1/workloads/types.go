package workloads

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

type Response struct {
	Meta           shared.PageMeta       `json:"meta"`
	Pods           Metric                `json:"pods"`
	Deployments    Metric                `json:"deployments"`
	ReplicaSets    Metric                `json:"replicaSets"`
	DaemonSets     Metric                `json:"daemonSets"`
	StatefulSets   Metric                `json:"statefulSets"`
	CronJobs       Metric                `json:"cronJobs"`
	Jobs           Metric                `json:"jobs"`
	ResourceQuotas Metric                `json:"resourceQuotas"`
	Warnings       []shared.WarningEvent `json:"warnings"`
}
