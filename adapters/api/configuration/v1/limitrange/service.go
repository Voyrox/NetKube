package limitrange

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	limitRanges, err := clientset.CoreV1().LimitRanges("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(limitRanges.Items))
	for _, item := range limitRanges.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Limits: len(item.Spec.Limits), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}
