package pod

import (
	"net/http"
	"strings"

	"netkube/adapters/api/shared"

	"github.com/gin-gonic/gin"
)

func ListHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := shared.SelectedNamespace(c)
	items, stats, err := List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Items: items, Count: len(items), Stats: stats})
}

func CreateHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	var request shared.ManifestCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "invalid pod manifest request"})
		return
	}

	pod, err := Create(cluster.Clientset, request.Content)
	if err != nil {
		c.JSON(shared.CreateStatusCode(err), shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shared.CreatedResourceResponse{Meta: shared.PageMetaFromCluster(cluster, pod.Namespace), Name: pod.Name, Namespace: pod.Namespace, Kind: "Pod"})
}

func DetailHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "pod name and namespace are required"})
		return
	}

	item, err := DetailFor(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, DetailResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Item: item})
}

func LogsHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	container := strings.TrimSpace(c.Query("container"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "pod name and namespace are required"})
		return
	}

	content, selectedContainer, err := Logs(cluster.Clientset, namespace, name, container)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, LogResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Container: selectedContainer, Content: content})
}

func EventsHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "pod name and namespace are required"})
		return
	}

	items, err := Events(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, EventsResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Items: items})
}

func YAMLHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "pod name and namespace are required"})
		return
	}

	content, err := YAML(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, YAMLResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Content: content})
}
