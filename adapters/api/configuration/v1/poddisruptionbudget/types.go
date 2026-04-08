package poddisruptionbudget

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace      string `json:"namespace"`
	Name           string `json:"name"`
	MinAvailable   string `json:"minAvailable"`
	MaxUnavailable string `json:"maxUnavailable"`
	Allowed        int32  `json:"allowed"`
	Age            string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
