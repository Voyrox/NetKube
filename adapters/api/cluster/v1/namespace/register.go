package namespace

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
	items, err := List(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ListResponse{Meta: shared.PageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func YAMLHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "namespace name is required"})
		return
	}
	content, err := YAML(cluster.Clientset, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, YAMLResponse{Meta: shared.PageMetaFromCluster(cluster, ""), Content: content})
}
