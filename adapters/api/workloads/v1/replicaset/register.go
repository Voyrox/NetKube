package replicaset

import (
	"net/http"

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
