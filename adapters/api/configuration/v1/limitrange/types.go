package limitrange

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Limits    int    `json:"limits"`
	Age       string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
