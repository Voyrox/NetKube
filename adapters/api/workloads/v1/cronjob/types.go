package cronjob

import "netkube/adapters/api/shared"

type Row struct {
	Namespace    string `json:"namespace"`
	Name         string `json:"name"`
	Schedule     string `json:"schedule"`
	Status       string `json:"status"`
	Suspend      string `json:"suspend"`
	Active       int    `json:"active"`
	LastSchedule string `json:"lastSchedule"`
	Age          string `json:"age"`
}

type Stats struct {
	Scheduled int `json:"scheduled"`
	Suspended int `json:"suspended"`
	Active    int `json:"active"`
	Total     int `json:"total"`
}

type ListResponse struct {
	Meta  shared.PageMeta `json:"meta"`
	Items []Row           `json:"items"`
	Count int             `json:"count"`
	Stats Stats           `json:"stats"`
}
