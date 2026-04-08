package resourcequota

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netkube/adapters/api/shared"
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
