package statefulset

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
	items, err := clientset.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
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

		desired := shared.DesiredReplicas(item.Spec.Replicas)
		rows = append(rows, Row{Namespace: item.Namespace, Name: item.Name, Ready: shared.ReadyRatio(item.Status.ReadyReplicas, desired), Status: resourceStatus, Desired: desired, Current: item.Status.CurrentReplicas, Updated: item.Status.UpdatedReplicas, Age: shared.FormatAge(item.CreationTimestamp)})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Namespace == rows[j].Namespace {
			return strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
		}
		return strings.ToLower(rows[i].Namespace) < strings.ToLower(rows[j].Namespace)
	})

	return rows, stats, nil
}

func status(item appsv1.StatefulSet) string {
	desired := shared.DesiredReplicas(item.Spec.Replicas)
	switch {
	case desired == 0 && item.Status.Replicas == 0:
		return "Scaled down"
	case item.Status.ReadyReplicas == desired && item.Status.UpdatedReplicas == desired:
		return "Healthy"
	case item.Status.ReadyReplicas == 0:
		return "Pending"
	case item.Status.UpdatedReplicas < item.Status.Replicas || item.Status.ReadyReplicas < desired:
		return "Updating"
	default:
		return "Degraded"
	}
}
