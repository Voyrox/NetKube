package volumeattributeclass

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	classes, err := clientset.StorageV1().VolumeAttributesClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(classes.Items))
	for _, item := range classes.Items {
		items = append(items, ListItem{Name: item.Name, DriverName: shared.Fallback(item.DriverName), Parameters: len(item.Parameters), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}
