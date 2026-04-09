package deployment

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
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "invalid deployment manifest request"})
		return
	}

	deployment, err := Create(cluster.Clientset, request.Content)
	if err != nil {
		c.JSON(shared.CreateStatusCode(err), shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shared.CreatedResourceResponse{Meta: shared.PageMetaFromCluster(cluster, deployment.Namespace), Name: deployment.Name, Namespace: deployment.Namespace, Kind: "Deployment"})
}

func DeleteHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "deployment name and namespace are required"})
		return
	}

	err := Delete(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(shared.DeleteStatusCode(err), shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, shared.DeletedResourceResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Name: name, Namespace: namespace, Kind: "Deployment"})
}

func DetailHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "deployment name and namespace are required"})
		return
	}

	item, err := DetailFor(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, DetailResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Item: item})
}

func EventsHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "deployment name and namespace are required"})
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
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "deployment name and namespace are required"})
		return
	}

	content, err := YAML(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, YAMLResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Content: content})
}
