package hpa

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Target    string `json:"target"`
	Min       string `json:"min"`
	Max       int32  `json:"max"`
	Current   string `json:"current"`
	Age       string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
