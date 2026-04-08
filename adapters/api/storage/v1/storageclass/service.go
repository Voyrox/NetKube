package storageclass

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	classes, err := clientset.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(classes.Items))
	for _, item := range classes.Items {
		items = append(items, ListItem{Name: item.Name, Provisioner: shared.Fallback(item.Provisioner), ReclaimPolicy: reclaimPolicy(item.ReclaimPolicy), BindingMode: bindingMode(item.VolumeBindingMode), Default: shared.YesNo(isDefault(item)), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}
func reclaimPolicy(value *corev1.PersistentVolumeReclaimPolicy) string {
	if value == nil {
		return "-"
	}
	return string(*value)
}
func bindingMode(value *storagev1.VolumeBindingMode) string {
	if value == nil {
		return "-"
	}
	return string(*value)
}
func isDefault(item storagev1.StorageClass) bool {
	return item.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" || item.Annotations["storageclass.beta.kubernetes.io/is-default-class"] == "true"
}
