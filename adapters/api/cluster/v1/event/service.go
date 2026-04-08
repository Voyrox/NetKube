package event

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DetailFor(clientset *kubernetes.Clientset, namespace, name, reason, kind string) (Detail, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return Detail{}, err
	}
	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].CreationTimestamp.Time.After(events.Items[j].CreationTimestamp.Time)
	})
	if len(events.Items) == 0 {
		return Detail{}, nil
	}
	selected := events.Items[0]
	for _, item := range events.Items {
		if (namespace == "" || item.Namespace == namespace) && (name == "" || item.InvolvedObject.Name == name) && (reason == "" || item.Reason == reason) && (kind == "" || item.InvolvedObject.Kind == kind) {
			selected = item
			break
		}
	}

	related := make([]TimelineRow, 0, 5)
	for _, item := range events.Items {
		if item.InvolvedObject.Kind == selected.InvolvedObject.Kind && item.InvolvedObject.Name == selected.InvolvedObject.Name {
			related = append(related, TimelineRow{Title: shared.Fallback(item.Reason), Message: shared.Fallback(item.Message), Age: shared.FormatAge(item.CreationTimestamp), Type: shared.Fallback(item.Type)})
			if len(related) == 5 {
				break
			}
		}
	}

	return Detail{Title: shared.Fallback(selected.Message), Type: shared.Fallback(selected.Type), Namespace: shared.Fallback(selected.Namespace), Reason: shared.Fallback(selected.Reason), InvolvedObject: strings.TrimSpace(strings.Join([]string{shared.Fallback(selected.InvolvedObject.Kind), shared.Fallback(selected.InvolvedObject.Name)}, " / ")), Kind: shared.Fallback(selected.InvolvedObject.Kind), Name: shared.Fallback(selected.InvolvedObject.Name), Source: strings.TrimSpace(strings.Join([]string{shared.Fallback(selected.Source.Component), shared.Fallback(selected.Source.Host)}, " / ")), FirstSeen: shared.FormatAge(selected.CreationTimestamp), LastSeen: shared.FormatAge(selected.CreationTimestamp), Count: selected.Count, Node: shared.Fallback(selected.Source.Host), Message: shared.Fallback(selected.Message), Timeline: related, Annotations: shared.CloneStringMap(selected.Annotations)}, nil
}
