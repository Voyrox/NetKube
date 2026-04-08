package service

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	ExternalIP string `json:"externalIP"`
	Ports      string `json:"ports"`
	Selector   string `json:"selector"`
	Age        string `json:"age"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
