package csinode

import "netkube/adapters/api/shared"

type ListItem struct {
	Name    string `json:"name"`
	Drivers int    `json:"drivers"`
	Age     string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
