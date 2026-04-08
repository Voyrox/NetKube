package secret

import (
	"context"
	"encoding/base64"
	"sort"
	"unicode/utf8"

	"netkube/adapters/api/shared"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	secrets, err := clientset.CoreV1().Secrets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(secrets.Items))
	for _, item := range secrets.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Type: shared.Fallback(string(item.Type)), Data: len(item.Data), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}
func Data(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	content := make(map[string]string, len(secret.Data))
	for key, value := range secret.Data {
		if utf8.Valid(value) {
			content[key] = string(value)
		} else {
			content[key] = "base64:" + base64.StdEncoding.EncodeToString(value)
		}
	}
	rendered, err := yaml.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(rendered), nil
}
