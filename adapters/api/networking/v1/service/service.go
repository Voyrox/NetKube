package service

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	services, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]ListItem, 0, len(services.Items))
	for _, item := range services.Items {
		items = append(items, ListItem{Namespace: item.Namespace, Name: item.Name, Type: shared.Fallback(string(item.Spec.Type)), ExternalIP: externalIP(item), Ports: ports(item), Selector: shared.MapToSelector(item.Spec.Selector), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}

func externalIP(service corev1.Service) string {
	if len(service.Spec.ExternalIPs) > 0 {
		return strings.Join(service.Spec.ExternalIPs, ", ")
	}
	if len(service.Status.LoadBalancer.Ingress) > 0 {
		values := make([]string, 0, len(service.Status.LoadBalancer.Ingress))
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			if ingress.IP != "" {
				values = append(values, ingress.IP)
				continue
			}
			if ingress.Hostname != "" {
				values = append(values, ingress.Hostname)
			}
		}
		if len(values) > 0 {
			return strings.Join(values, ", ")
		}
	}
	if service.Spec.ClusterIP != "" && service.Spec.ClusterIP != "None" {
		return service.Spec.ClusterIP
	}
	return "-"
}

func ports(service corev1.Service) string {
	if len(service.Spec.Ports) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(service.Spec.Ports))
	for _, port := range service.Spec.Ports {
		part := string(port.Protocol) + "/" + shared.Int32String(port.Port)
		if port.TargetPort.String() != "" {
			part += " -> " + port.TargetPort.String()
		}
		if strings.TrimSpace(port.Name) != "" {
			part += " (" + port.Name + ")"
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, ", ")
}
