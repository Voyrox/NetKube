package lease

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	coordinationv1 "k8s.io/api/coordination/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	leases, err := clientset.CoordinationV1().Leases("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(leases.Items))
	for _, item := range leases.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Holder: holder(item), LastRenew: renew(item), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}

func YAML(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	lease, err := clientset.CoordinationV1().Leases(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	content, err := yaml.Marshal(lease)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func holder(lease coordinationv1.Lease) string {
	if lease.Spec.HolderIdentity == nil || strings.TrimSpace(*lease.Spec.HolderIdentity) == "" {
		return "-"
	}
	return *lease.Spec.HolderIdentity
}

func renew(lease coordinationv1.Lease) string {
	if lease.Spec.RenewTime == nil || lease.Spec.RenewTime.IsZero() {
		return "-"
	}
	return shared.FormatAge(metav1.Time{Time: lease.Spec.RenewTime.Time})
}
