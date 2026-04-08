package storageclass

import "netkube/adapters/api/shared"

type ListItem struct {
	Name          string `json:"name"`
	Provisioner   string `json:"provisioner"`
	ReclaimPolicy string `json:"reclaimPolicy"`
	BindingMode   string `json:"bindingMode"`
	Default       string `json:"default"`
	Age           string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
