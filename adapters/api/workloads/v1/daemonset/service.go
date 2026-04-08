package daemonset

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func List(clientset *kubernetes.Clientset, namespace string) ([]Row, Stats, error) {
	items, err := clientset.AppsV1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, Stats{}, err
	}

	rows := make([]Row, 0, len(items.Items))
	stats := Stats{Total: len(items.Items)}
	for _, item := range items.Items {
		resourceStatus := status(item)
		switch resourceStatus {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		rows = append(rows, Row{Namespace: item.Namespace, Name: item.Name, Ready: shared.ReadyRatio(item.Status.NumberReady, item.Status.DesiredNumberScheduled), Status: resourceStatus, Desired: item.Status.DesiredNumberScheduled, Current: item.Status.CurrentNumberScheduled, Available: item.Status.NumberAvailable, Misscheduled: item.Status.NumberMisscheduled, Age: shared.FormatAge(item.CreationTimestamp)})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Namespace == rows[j].Namespace {
			return strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
		}
		return strings.ToLower(rows[i].Namespace) < strings.ToLower(rows[j].Namespace)
	})

	return rows, stats, nil
}

func status(item appsv1.DaemonSet) string {
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
