package api

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type deploymentRow struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Ready     string `json:"ready"`
	Status    string `json:"status"`
	Desired   int32  `json:"desired"`
	Updated   int32  `json:"updated"`
	Available int32  `json:"available"`
	Age       string `json:"age"`
}

type deploymentsResponse struct {
	Meta  pageMeta         `json:"meta"`
	Items []deploymentRow  `json:"items"`
	Count int              `json:"count"`
	Error string           `json:"error,omitempty"`
	Stats deploymentsStats `json:"stats"`
}

type deploymentsStats struct {
	Healthy int `json:"healthy"`
	Warning int `json:"warning"`
	Pending int `json:"pending"`
	Total   int `json:"total"`
}

func DeploymentsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetDeployments(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deploymentsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func GetDeployments(clientset *kubernetes.Clientset, namespace string) ([]deploymentRow, deploymentsStats, error) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, deploymentsStats{}, err
	}

	items := make([]deploymentRow, 0, len(deployments.Items))
	stats := deploymentsStats{Total: len(deployments.Items)}

	for _, item := range deployments.Items {
		status := deploymentStatus(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		items = append(items, deploymentRow{
			Namespace: item.Namespace,
			Name:      item.Name,
			Ready:     readyRatio(item.Status.ReadyReplicas, item.Status.Replicas),
			Status:    status,
			Desired:   item.Status.Replicas,
			Updated:   item.Status.UpdatedReplicas,
			Available: item.Status.AvailableReplicas,
			Age:       formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		}

		return strings.ToLower(items[i].Namespace) < strings.ToLower(items[j].Namespace)
	})

	return items, stats, nil
}

func deploymentStatus(item appsv1.Deployment) string {
	switch {
	case item.Status.Replicas == 0:
		return "Pending"
	case item.Status.AvailableReplicas == item.Status.Replicas && item.Status.UpdatedReplicas == item.Status.Replicas:
		return "Healthy"
	case item.Status.UnavailableReplicas > 0:
		return "Degraded"
	default:
		return "Updating"
	}
}
