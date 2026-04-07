package storage

type SourceItem struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Size         int64  `json:"size,omitempty"`
	LastModified int64  `json:"lastModified,omitempty"`
	StoredFile   string `json:"storedFile,omitempty"`
	Content      string `json:"content,omitempty"`
}

type SourcesPayload struct {
	Sources []SourceItem `json:"sources"`
}

type SelectedContextsPayload struct {
	SelectedContextIDs []string `json:"selectedContextIds"`
}

type ContextItem struct {
	ID          string `json:"id"`
	ContextName string `json:"contextName"`
	ClusterName string `json:"clusterName"`
	UserName    string `json:"userName"`
	Namespace   string `json:"namespace"`
	SourceID    string `json:"sourceId"`
	SourceName  string `json:"sourceName"`
	Server      string `json:"server"`
	IsCurrent   bool   `json:"isCurrent"`
	HasUser     bool   `json:"hasUser"`
}

type ContextsPayload struct {
	Contexts []ContextItem `json:"contexts"`
}
