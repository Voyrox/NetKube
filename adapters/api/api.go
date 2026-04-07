package api

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"netkube/adapters"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const contextHeader = "X-NetKube-Context"

type apiError struct {
	Error string `json:"error"`
}

type pageMeta struct {
	ContextName string `json:"contextName"`
	ClusterName string `json:"clusterName"`
	UserName    string `json:"userName"`
	Namespace   string `json:"namespace,omitempty"`
	LastRefresh string `json:"lastRefresh"`
}

type warningEvent struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Age       string `json:"age"`
}

func resolveClusterRequest(c *gin.Context) (*adapters.ClusterClient, bool) {
	contextID := strings.TrimSpace(c.GetHeader(contextHeader))
	if contextID == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "missing X-NetKube-Context header"})
		return nil, false
	}

	cluster, err := adapters.ResolveCluster(contextID)
	if err != nil {
		status := http.StatusBadGateway
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "missing") {
			status = http.StatusNotFound
		}

		c.JSON(status, apiError{Error: err.Error()})
		return nil, false
	}

	return cluster, true
}

func selectedNamespace(c *gin.Context) string {
	namespace := strings.TrimSpace(c.Query("namespace"))
	if namespace == "" || strings.EqualFold(namespace, "all") {
		return ""
	}

	return namespace
}

func pageMetaFromCluster(cluster *adapters.ClusterClient, namespace string) pageMeta {
	return pageMeta{
		ContextName: cluster.Context.ContextName,
		ClusterName: cluster.Context.ClusterName,
		UserName:    cluster.Context.UserName,
		Namespace:   namespace,
		LastRefresh: time.Now().Format(time.RFC3339),
	}
}

func listWarningEvents(clientset *kubernetes.Clientset, namespace string, limit int) []warningEvent {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return []warningEvent{}
	}

	filtered := make([]warningEvent, 0, limit)
	sort.Slice(events.Items, func(i, j int) bool {
		return eventTimestamp(&events.Items[i].ObjectMeta).After(eventTimestamp(&events.Items[j].ObjectMeta))
	})

	for _, item := range events.Items {
		if item.Type != "Warning" {
			continue
		}

		filtered = append(filtered, warningEvent{
			Namespace: item.Namespace,
			Name:      item.InvolvedObject.Name,
			Reason:    item.Reason,
			Message:   item.Message,
			Age:       formatAge(item.CreationTimestamp),
		})

		if len(filtered) == limit {
			break
		}
	}

	return filtered
}

func eventTimestamp(event *metav1.ObjectMeta) time.Time {
	if event == nil {
		return time.Time{}
	}

	return event.CreationTimestamp.Time
}

func formatAge(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "-"
	}

	return formatDuration(time.Since(timestamp.Time))
}

func formatDuration(duration time.Duration) string {
	if duration < time.Minute {
		seconds := int(duration.Seconds())
		if seconds < 0 {
			seconds = 0
		}
		return fmt.Sprintf("%ds", seconds)
	}
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}
	if duration < 24*time.Hour {
		return fmt.Sprintf("%dh", int(duration.Hours()))
	}

	return fmt.Sprintf("%dd", int(duration.Hours()/24))
}
