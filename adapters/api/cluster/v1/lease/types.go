package lease

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Holder    string `json:"holder"`
	LastRenew string `json:"lastRenew"`
	Age       string `json:"age"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}

type YAMLResponse struct {
	Meta    shared.PageMeta `json:"meta"`
	Content string          `json:"content"`
}
