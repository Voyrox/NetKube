package csidriver

import (
	"context"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
	"strings"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	drivers, err := clientset.StorageV1().CSIDrivers().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(drivers.Items))
	for _, item := range drivers.Items {
		items = append(items, ListItem{Name: item.Name, AttachRequired: shared.YesNo(shared.BoolPointer(item.Spec.AttachRequired)), PodInfoOnMount: shared.YesNo(shared.BoolPointer(item.Spec.PodInfoOnMount)), StorageCapacity: shared.YesNo(shared.BoolPointer(item.Spec.StorageCapacity)), Modes: lifecycleModes(item.Spec.VolumeLifecycleModes), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}
func lifecycleModes(modes []storagev1.VolumeLifecycleMode) string {
	if len(modes) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(modes))
	for _, mode := range modes {
		parts = append(parts, string(mode))
	}
	return strings.Join(parts, ", ")
}
