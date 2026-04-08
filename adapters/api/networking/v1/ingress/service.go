package ingress

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(ingresses.Items))
	for _, item := range ingresses.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Class: class(item), Hosts: hosts(item), Address: address(item), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}

func class(ingress networkingv1.Ingress) string {
	if ingress.Spec.IngressClassName != nil && strings.TrimSpace(*ingress.Spec.IngressClassName) != "" {
		return *ingress.Spec.IngressClassName
	}
	return "-"
}

func hosts(ingress networkingv1.Ingress) string {
	if len(ingress.Spec.Rules) == 0 {
		return "-"
	}
	hosts := make([]string, 0, len(ingress.Spec.Rules))
	for _, rule := range ingress.Spec.Rules {
		if strings.TrimSpace(rule.Host) != "" {
			hosts = append(hosts, rule.Host)
		}
	}
	if len(hosts) == 0 {
		return "-"
	}
	return strings.Join(hosts, ", ")
}

func address(ingress networkingv1.Ingress) string {
	if len(ingress.Status.LoadBalancer.Ingress) == 0 {
		return "-"
	}
	addresses := make([]string, 0, len(ingress.Status.LoadBalancer.Ingress))
	for _, entry := range ingress.Status.LoadBalancer.Ingress {
		if strings.TrimSpace(entry.IP) != "" {
			addresses = append(addresses, entry.IP)
			continue
		}
		if strings.TrimSpace(entry.Hostname) != "" {
			addresses = append(addresses, entry.Hostname)
		}
	}
	if len(addresses) == 0 {
		return "-"
	}
	return strings.Join(addresses, ", ")
}
