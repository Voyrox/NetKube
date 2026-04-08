package csidriver

import "netkube/adapters/api/shared"

type ListItem struct {
	Name            string `json:"name"`
	AttachRequired  string `json:"attachRequired"`
	PodInfoOnMount  string `json:"podInfoOnMount"`
	StorageCapacity string `json:"storageCapacity"`
	Modes           string `json:"modes"`
	Age             string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
