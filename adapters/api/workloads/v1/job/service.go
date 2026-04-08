package job

import (
	"context"
	"sort"
	"strings"
	"time"

	"netkube/adapters/api/shared"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func List(clientset *kubernetes.Clientset, namespace string) ([]Row, Stats, error) {
	items, err := clientset.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, Stats{}, err
	}

	rows := make([]Row, 0, len(items.Items))
	stats := Stats{Total: len(items.Items)}
	for _, item := range items.Items {
		resourceStatus := status(item)
		switch resourceStatus {
		case "Succeeded":
			stats.Succeeded++
		case "Running":
			stats.Active++
		case "Failed":
			stats.Failed++
		}

		rows = append(rows, Row{Namespace: item.Namespace, Name: item.Name, Status: resourceStatus, Completions: shared.ReadyRatio(item.Status.Succeeded, desiredCompletions(item.Spec.Completions)), Active: item.Status.Active, Duration: duration(item), Age: shared.FormatAge(item.CreationTimestamp)})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Namespace == rows[j].Namespace {
			return strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
		}
		return strings.ToLower(rows[i].Namespace) < strings.ToLower(rows[j].Namespace)
	})

	return rows, stats, nil
}

func status(item batchv1.Job) string {
	switch {
	case item.Status.Failed > 0:
		return "Failed"
	case item.Status.Succeeded > 0:
		return "Succeeded"
	case item.Status.Active > 0:
		return "Running"
	default:
		return "Pending"
	}
}

func desiredCompletions(value *int32) int32 {
	if value == nil {
		return 1
	}
	return *value
}

func duration(item batchv1.Job) string {
	if item.Status.StartTime == nil {
		return "-"
	}
	if item.Status.CompletionTime != nil {
		return shared.FormatDuration(item.Status.CompletionTime.Sub(item.Status.StartTime.Time))
	}
	return shared.FormatDuration(time.Since(item.Status.StartTime.Time))
}
