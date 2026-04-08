package replicaset

import "netkube/adapters/api/shared"

type Row struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Ready     string `json:"ready"`
	Status    string `json:"status"`
	Desired   int32  `json:"desired"`
	Current   int32  `json:"current"`
	ReadyPods int32  `json:"readyPods"`
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
	Stats Stats           `json:"stats"`
}
