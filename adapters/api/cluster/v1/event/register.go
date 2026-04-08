package event

import (
	"net/http"
	"strings"

	"netkube/adapters/api/shared"

	"github.com/gin-gonic/gin"
)

func DetailHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}
	item, err := DetailFor(cluster.Clientset, strings.TrimSpace(c.Query("namespace")), strings.TrimSpace(c.Query("name")), strings.TrimSpace(c.Query("reason")), strings.TrimSpace(c.Query("kind")))
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, DetailResponse{Meta: shared.PageMetaFromCluster(cluster, item.Namespace), Item: item})
}
