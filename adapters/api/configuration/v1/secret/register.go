package secret

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netkube/adapters/api/shared"
	"strings"
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
func DataHandler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}
	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, shared.APIError{Error: "secret name and namespace are required"})
		return
	}
	content, err := Data(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, DataResponse{Meta: shared.PageMetaFromCluster(cluster, namespace), Namespace: namespace, Name: name, Content: content})
}
