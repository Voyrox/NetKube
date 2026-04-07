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

func ReplicaSetsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetReplicaSets(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, replicaSetsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func DaemonSetsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetDaemonSets(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, daemonSetsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func GetReplicaSets(clientset *kubernetes.Clientset, namespace string) ([]replicaSetRow, replicaSetsStats, error) {
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, replicaSetsStats{}, err
	}

	items := make([]replicaSetRow, 0, len(replicaSets.Items))
	stats := replicaSetsStats{Total: len(replicaSets.Items)}

	for _, item := range replicaSets.Items {
		status := replicaSetStatus(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		desired := replicasOrZero(item.Spec.Replicas)
		items = append(items, replicaSetRow{
			Namespace: item.Namespace,
			Name:      item.Name,
			Ready:     readyRatio(item.Status.ReadyReplicas, desired),
			Status:    status,
			Desired:   desired,
			Current:   item.Status.Replicas,
			ReadyPods: item.Status.ReadyReplicas,
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

func GetDaemonSets(clientset *kubernetes.Clientset, namespace string) ([]daemonSetRow, daemonSetsStats, error) {
	daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, daemonSetsStats{}, err
	}

	items := make([]daemonSetRow, 0, len(daemonSets.Items))
	stats := daemonSetsStats{Total: len(daemonSets.Items)}

	for _, item := range daemonSets.Items {
		status := daemonSetStatus(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		items = append(items, daemonSetRow{
			Namespace:    item.Namespace,
			Name:         item.Name,
			Ready:        readyRatio(item.Status.NumberReady, item.Status.DesiredNumberScheduled),
			Status:       status,
			Desired:      item.Status.DesiredNumberScheduled,
			Current:      item.Status.CurrentNumberScheduled,
			Available:    item.Status.NumberAvailable,
			Misscheduled: item.Status.NumberMisscheduled,
			Age:          formatAge(item.CreationTimestamp),
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

func replicaSetStatus(item appsv1.ReplicaSet) string {
	desired := replicasOrZero(item.Spec.Replicas)
	switch {
	case desired == 0 && item.Status.Replicas == 0:
		return "Scaled down"
	case desired == 0:
		return "Updating"
	case item.Status.ReadyReplicas == desired && item.Status.AvailableReplicas == desired:
		return "Healthy"
	case item.Status.ReadyReplicas == 0:
		return "Degraded"
	default:
		return "Updating"
	}
}

func daemonSetStatus(item appsv1.DaemonSet) string {
	desired := item.Status.DesiredNumberScheduled
	switch {
	case desired == 0 && item.Status.CurrentNumberScheduled == 0:
		return "Scaled down"
	case desired == 0:
		return "Updating"
	case item.Status.NumberReady == desired && item.Status.UpdatedNumberScheduled == desired && item.Status.NumberMisscheduled == 0:
		return "Healthy"
	case item.Status.NumberReady == 0 || item.Status.NumberMisscheduled > 0:
		return "Degraded"
	default:
		return "Updating"
	}
}

func replicasOrZero(value *int32) int32 {
	if value == nil {
		return 0
	}

	return *value
}
