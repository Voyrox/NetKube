package job

import "netkube/adapters/api/shared"

type Row struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Completions string `json:"completions"`
	Active      int32  `json:"active"`
	Duration    string `json:"duration"`
	Age         string `json:"age"`
}

type Stats struct {
	Succeeded int `json:"succeeded"`
	Active    int `json:"active"`
	Failed    int `json:"failed"`
	Total     int `json:"total"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []Row           `json:"items"`
	Count int             `json:"count"`
	Stats Stats           `json:"stats"`
}
