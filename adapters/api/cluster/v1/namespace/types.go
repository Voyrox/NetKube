package namespace

import "netkube/adapters/api/shared"

type ListItem struct {
	Name  string `json:"name"`
	Phase string `json:"phase"`
	Age   string `json:"age"`
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
