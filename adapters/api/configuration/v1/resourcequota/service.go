package resourcequota

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	resourceQuotas, err := clientset.CoreV1().ResourceQuotas("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(resourceQuotas.Items))
	for _, item := range resourceQuotas.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Scopes: len(item.Spec.Scopes), Hard: len(item.Status.Hard), Used: len(item.Status.Used), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}
