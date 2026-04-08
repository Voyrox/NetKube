package volumeattachment

import "netkube/adapters/api/shared"

type ListItem struct {
	Name             string `json:"name"`
	Attacher         string `json:"attacher"`
	Node             string `json:"node"`
	PersistentVolume string `json:"persistentVolume"`
	Attached         string `json:"attached"`
	Age              string `json:"age"`
}
type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []ListItem      `json:"items"`
	Count int             `json:"count"`
}
