package hpa

import (
	"context"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
	"strings"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	hpas, err := clientset.AutoscalingV2().HorizontalPodAutoscalers("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(hpas.Items))
	for _, item := range hpas.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Target: target(item), Min: shared.Int32PointerString(item.Spec.MinReplicas), Max: item.Spec.MaxReplicas, Current: shared.Int32String(item.Status.CurrentReplicas), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}
func target(item autoscalingv2.HorizontalPodAutoscaler) string {
	kind := strings.TrimSpace(item.Spec.ScaleTargetRef.Kind)
	name := strings.TrimSpace(item.Spec.ScaleTargetRef.Name)
	if kind == "" && name == "" {
		return "-"
	}
	return strings.TrimSpace(strings.Join([]string{kind, name}, "/"))
}
