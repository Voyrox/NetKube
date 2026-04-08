package secret

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Data      int    `json:"data"`
	Age       string `json:"age"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
type DataResponse struct {
	Meta      shared.PageMeta `json:"meta"`
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Content   string          `json:"content"`
}
