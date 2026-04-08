package persistentvolume

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
	"strings"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	volumes, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(volumes.Items))
	for _, item := range volumes.Items {
		items = append(items, ListItem{Name: item.Name, Status: shared.Fallback(string(item.Status.Phase)), Capacity: resourceListValue(item.Spec.Capacity, corev1.ResourceStorage), AccessModes: accessModes(item.Spec.AccessModes), ReclaimPolicy: shared.Fallback(string(item.Spec.PersistentVolumeReclaimPolicy)), Claim: claim(item.Spec.ClaimRef), StorageClass: shared.Fallback(item.Spec.StorageClassName), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}
func accessModes(modes []corev1.PersistentVolumeAccessMode) string {
	if len(modes) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(modes))
	for _, mode := range modes {
		parts = append(parts, string(mode))
	}
	return strings.Join(parts, ", ")
}
func claim(ref *corev1.ObjectReference) string {
	if ref == nil {
		return "-"
	}
	if strings.TrimSpace(ref.Namespace) == "" {
		return shared.Fallback(ref.Name)
	}
	return strings.TrimSpace(ref.Namespace + "/" + ref.Name)
}
func resourceListValue(values corev1.ResourceList, key corev1.ResourceName) string {
	quantity, ok := values[key]
	if !ok {
		return "-"
	}
	return quantity.String()
}
