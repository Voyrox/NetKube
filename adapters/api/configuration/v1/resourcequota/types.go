package resourcequota

import "netkube/adapters/api/shared"

type ListItem struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Scopes    int    `json:"scopes"`
	Hard      int    `json:"hard"`
	Used      int    `json:"used"`
	Age       string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
