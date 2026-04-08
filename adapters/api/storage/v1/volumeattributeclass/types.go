package volumeattributeclass

import "netkube/adapters/api/shared"

type ListItem struct {
	Name       string `json:"name"`
	DriverName string `json:"driverName"`
	Parameters int    `json:"parameters"`
	Age        string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
