package persistentvolumeclaim

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
	claims, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(claims.Items))
	for _, item := range claims.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Status: shared.Fallback(string(item.Status.Phase)), Volume: shared.Fallback(item.Spec.VolumeName), Capacity: resourceListValue(item.Status.Capacity, corev1.ResourceStorage), AccessModes: pvcAccessModes(item.Status.AccessModes, item.Spec.AccessModes), StorageClass: pvcStorageClass(item.Spec.StorageClassName), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}
func pvcAccessModes(primary, fallbackModes []corev1.PersistentVolumeAccessMode) string {
	if len(primary) > 0 {
		return pvAccessModes(primary)
	}
	return pvAccessModes(fallbackModes)
}
func pvAccessModes(modes []corev1.PersistentVolumeAccessMode) string {
	if len(modes) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(modes))
	for _, mode := range modes {
		parts = append(parts, string(mode))
	}
	return strings.Join(parts, ", ")
}
func pvcStorageClass(name *string) string {
	if name == nil || strings.TrimSpace(*name) == "" {
		return "-"
	}
	return *name
}
func resourceListValue(values corev1.ResourceList, key corev1.ResourceName) string {
	quantity, ok := values[key]
	if !ok {
		return "-"
	}
	return quantity.String()
}
