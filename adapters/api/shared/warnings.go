package shared

import (
	"context"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListWarningEvents(clientset *kubernetes.Clientset, namespace string, limit int) []WarningEvent {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return []WarningEvent{}
	}

	filtered := make([]WarningEvent, 0, limit)
	sort.Slice(events.Items, func(i, j int) bool {
		return EventTimestamp(&events.Items[i].ObjectMeta).After(EventTimestamp(&events.Items[j].ObjectMeta))
	})

	for _, item := range events.Items {
		if item.Type != "Warning" {
			continue
		}

		filtered = append(filtered, WarningEvent{
			Namespace: item.Namespace,
			Name:      item.InvolvedObject.Name,
			Reason:    item.Reason,
			Message:   item.Message,
			Age:       FormatAge(item.CreationTimestamp),
		})

		if len(filtered) == limit {
			break
		}
	}

	return filtered
}
