package cronjob

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func List(clientset *kubernetes.Clientset, namespace string) ([]Row, Stats, error) {
	items, err := clientset.BatchV1().CronJobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, Stats{}, err
	}

	rows := make([]Row, 0, len(items.Items))
	stats := Stats{Total: len(items.Items)}
	for _, item := range items.Items {
		resourceStatus := status(item)
		if item.Spec.Suspend != nil && *item.Spec.Suspend {
			stats.Suspended++
		} else {
			stats.Scheduled++
		}
		if len(item.Status.Active) > 0 {
			stats.Active++
		}

		rows = append(rows, Row{Namespace: item.Namespace, Name: item.Name, Schedule: item.Spec.Schedule, Status: resourceStatus, Suspend: shared.YesNo(item.Spec.Suspend != nil && *item.Spec.Suspend), Active: len(item.Status.Active), LastSchedule: shared.FormatOptionalAge(item.Status.LastScheduleTime), Age: shared.FormatAge(item.CreationTimestamp)})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Namespace == rows[j].Namespace {
			return strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
		}
		return strings.ToLower(rows[i].Namespace) < strings.ToLower(rows[j].Namespace)
	})

	return rows, stats, nil
}

func status(item batchv1.CronJob) string {
	switch {
	case item.Spec.Suspend != nil && *item.Spec.Suspend:
		return "Suspended"
	case len(item.Status.Active) > 0:
		return "Running"
	case item.Status.LastScheduleTime == nil:
		return "Pending"
	default:
		return "Scheduled"
	}
}
