package api

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type podRow struct {
	Namespace         string `json:"namespace"`
	Name              string `json:"name"`
	Ready             string `json:"ready"`
	Status            string `json:"status"`
	Restarts          int32  `json:"restarts"`
	LastRestart       string `json:"lastRestart"`
	LastRestartReason string `json:"lastRestartReason"`
	Node              string `json:"node"`
	PodIP             string `json:"podIP"`
	Age               string `json:"age"`
}

type podsStats struct {
	Running int `json:"running"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
	Other   int `json:"other"`
	Total   int `json:"total"`
}

type podsResponse struct {
	Meta  pageMeta  `json:"meta"`
	Items []podRow  `json:"items"`
	Count int       `json:"count"`
	Stats podsStats `json:"stats"`
}

func PodsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetPods(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, podsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func GetPods(clientset *kubernetes.Clientset, namespace string) ([]podRow, podsStats, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, podsStats{}, err
	}

	items := make([]podRow, 0, len(pods.Items))
	stats := podsStats{Total: len(pods.Items)}

	for _, item := range pods.Items {
		status := podStatus(item)
		switch status {
		case "Running":
			stats.Running++
		case "Pending":
			stats.Pending++
		case "Failed":
			stats.Failed++
		default:
			stats.Other++
		}

		restarts, lastRestart, lastRestartReason := containerRestartInfo(item)
		items = append(items, podRow{
			Namespace:         item.Namespace,
			Name:              item.Name,
			Ready:             readyRatio(readyContainers(item), int32(len(item.Status.ContainerStatuses))),
			Status:            status,
			Restarts:          restarts,
			LastRestart:       lastRestart,
			LastRestartReason: lastRestartReason,
			Node:              fallback(item.Spec.NodeName),
			PodIP:             fallback(item.Status.PodIP),
			Age:               formatAge(item.CreationTimestamp),
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

func podStatus(pod corev1.Pod) string {
	for _, status := range pod.Status.ContainerStatuses {
		if status.State.Waiting != nil && status.State.Waiting.Reason != "" {
			return status.State.Waiting.Reason
		}
		if status.State.Terminated != nil && status.State.Terminated.Reason != "" {
			return status.State.Terminated.Reason
		}
	}

	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}

	if pod.Status.Phase == "" {
		return "Unknown"
	}

	return string(pod.Status.Phase)
}

func readyContainers(pod corev1.Pod) int32 {
	var ready int32
	for _, status := range pod.Status.ContainerStatuses {
		if status.Ready {
			ready++
		}
	}

	return ready
}

func containerRestartInfo(pod corev1.Pod) (int32, string, string) {
	var restarts int32
	lastRestart := "-"
	lastRestartReason := "-"
	var latest metav1.Time

	for _, status := range pod.Status.ContainerStatuses {
		restarts += status.RestartCount
		if status.LastTerminationState.Terminated == nil {
			continue
		}

		terminated := status.LastTerminationState.Terminated
		if latest.IsZero() || terminated.FinishedAt.After(latest.Time) {
			latest = terminated.FinishedAt
			lastRestart = formatAge(terminated.FinishedAt)
			lastRestartReason = fallback(terminated.Reason)
		}
	}

	return restarts, lastRestart, lastRestartReason
}

func readyRatio(current, total int32) string {
	return strings.TrimSpace(strings.Join([]string{int32String(current), int32String(total)}, "/"))
}

func int32String(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func fallback(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}

	return value
}
