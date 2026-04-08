package namespace

import (
	"context"
	"sort"

	"netkube/adapters/api/shared"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	items, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	rows := make([]ListItem, 0, len(items.Items))
	for _, item := range items.Items {
		rows = append(rows, ListItem{Name: item.Name, Phase: shared.Fallback(string(item.Status.Phase)), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows, nil
}

func YAML(clientset *kubernetes.Clientset, name string) (string, error) {
	resource, err := clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	content, err := yaml.Marshal(resource)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
