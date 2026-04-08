package configmap

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Data      int    `json:"data"`
	Immutable string `json:"immutable"`
	Age       string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
