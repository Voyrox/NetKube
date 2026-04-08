package volumeattachment

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"netkube/adapters/api/shared"
	"sort"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	attachments, err := clientset.StorageV1().VolumeAttachments().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(attachments.Items))
	for _, item := range attachments.Items {
		items = append(items, ListItem{Name: item.Name, Attacher: shared.Fallback(item.Spec.Attacher), Node: shared.Fallback(item.Spec.NodeName), PersistentVolume: shared.Fallback(shared.StringPointer(item.Spec.Source.PersistentVolumeName)), Attached: shared.YesNo(item.Status.Attached), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}
