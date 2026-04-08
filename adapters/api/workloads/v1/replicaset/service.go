package replicaset

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
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, Stats{}, err
	}

	items := make([]Row, 0, len(replicaSets.Items))
	stats := Stats{Total: len(replicaSets.Items)}
	for _, item := range replicaSets.Items {
		resourceStatus := status(item)
		switch resourceStatus {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		desired := shared.ReplicasOrZero(item.Spec.Replicas)
		items = append(items, Row{Namespace: item.Namespace, Name: item.Name, Ready: shared.ReadyRatio(item.Status.ReadyReplicas, desired), Status: resourceStatus, Desired: desired, Current: item.Status.Replicas, ReadyPods: item.Status.ReadyReplicas, Age: shared.FormatAge(item.CreationTimestamp)})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		}
		return strings.ToLower(items[i].Namespace) < strings.ToLower(items[j].Namespace)
	})

	return items, stats, nil
}

func status(item appsv1.ReplicaSet) string {
	desired := shared.ReplicasOrZero(item.Spec.Replicas)
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
