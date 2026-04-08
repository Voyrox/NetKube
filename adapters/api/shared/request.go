package shared

import (
	"net/http"
	"strings"

	"netkube/adapters"

	"github.com/gin-gonic/gin"
)

const ContextHeader = "X-NetKube-Context"

func ResolveClusterRequest(c *gin.Context) (*adapters.ClusterClient, bool) {
	contextID := strings.TrimSpace(c.GetHeader(ContextHeader))
	if contextID == "" {
		c.JSON(http.StatusBadRequest, APIError{Error: "missing X-NetKube-Context header"})
		return nil, false
	}

	cluster, err := adapters.ResolveCluster(contextID)
	if err != nil {
		status := http.StatusBadGateway
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "missing") {
			status = http.StatusNotFound
		}

		c.JSON(status, APIError{Error: err.Error()})
		return nil, false
	}

	return cluster, true
}

func SelectedNamespace(c *gin.Context) string {
	namespace := strings.TrimSpace(c.Query("namespace"))
	if namespace == "" || strings.EqualFold(namespace, "all") {
		return ""
	}

	return namespace
}
