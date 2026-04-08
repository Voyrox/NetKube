package node

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func NodesSummary(clientset *kubernetes.Clientset) (Summary, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return Summary{}, err
	}

	summary := Summary{Total: len(nodes.Items)}
	for _, node := range nodes.Items {
		if ready(node) {
			summary.Ready++
			continue
		}
		if node.Spec.Unschedulable {
			summary.Pending++
			continue
		}
		summary.NotReady++
	}

	switch {
	case summary.Total == 0:
		summary.Status = "Unknown"
	case summary.NotReady > 0:
		summary.Status = "Attention"
	case summary.Pending > 0:
		summary.Status = "Watch"
	default:
		summary.Status = "Healthy"
	}

	return summary, nil
}

func DetailFor(clientset *kubernetes.Clientset, name string) (Detail, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return Detail{}, err
	}
	if len(nodes.Items) == 0 {
		return Detail{}, nil
	}

	sort.Slice(nodes.Items, func(i, j int) bool { return nodes.Items[i].Name < nodes.Items[j].Name })
	node := nodes.Items[0]
	if name != "" {
		for _, item := range nodes.Items {
			if item.Name == name {
				node = item
				break
			}
		}
	}

	conditions := make([]ConditionRow, 0, len(node.Status.Conditions))
	var cpuPressure, memoryPressure, diskPressure string
	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, ConditionRow{Type: string(condition.Type), Status: string(condition.Status), Reason: shared.Fallback(condition.Reason), Message: shared.Fallback(condition.Message)})
		switch condition.Type {
		case corev1.NodePIDPressure:
			if cpuPressure == "" {
				cpuPressure = boolStateLabel(condition.Status == corev1.ConditionFalse)
			}
		case corev1.NodeMemoryPressure:
			memoryPressure = boolStateLabel(condition.Status == corev1.ConditionFalse)
		case corev1.NodeDiskPressure:
			diskPressure = boolStateLabel(condition.Status == corev1.ConditionFalse)
		}
	}

	events, _ := clientset.CoreV1().Events("").List(context.Background(), metav1.ListOptions{FieldSelector: "involvedObject.kind=Node,involvedObject.name=" + node.Name})
	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].CreationTimestamp.Time.After(events.Items[j].CreationTimestamp.Time)
	})
	timeline := make([]EventRow, 0, shared.MinInt(len(events.Items), 5))
	for _, item := range events.Items {
		timeline = append(timeline, EventRow{Title: shared.Fallback(item.Reason), Message: shared.Fallback(item.Message), Age: shared.FormatAge(item.CreationTimestamp), Type: shared.Fallback(item.Type)})
		if len(timeline) == 5 {
			break
		}
	}

	allocatableStore := int64(0)
	if quantity, ok := node.Status.Allocatable[corev1.ResourceEphemeralStorage]; ok {
		allocatableStore = quantity.Value()
	}

	return Detail{Name: node.Name, Status: status(node), Role: role(node), KubeletVersion: shared.Fallback(node.Status.NodeInfo.KubeletVersion), ContainerRuntime: shared.Fallback(node.Status.NodeInfo.ContainerRuntimeVersion), OSKernel: strings.TrimSpace(strings.Join([]string{shared.Fallback(node.Status.NodeInfo.OSImage), shared.Fallback(node.Status.NodeInfo.KernelVersion)}, " / ")), Architecture: shared.Fallback(node.Status.NodeInfo.Architecture), InternalIP: address(node, corev1.NodeInternalIP), PodCIDR: shared.Fallback(node.Spec.PodCIDR), AllocatableCPU: shared.FormatCPU(node.Status.Allocatable.Cpu().MilliValue()), AllocatableMemory: shared.FormatBytes(node.Status.Allocatable.Memory().Value()), AllocatablePods: shared.Fallback(node.Status.Allocatable.Pods().String()), AllocatableStore: shared.FormatBytes(allocatableStore), CPUPressure: shared.Fallback(cpuPressure), MemoryPressure: shared.Fallback(memoryPressure), DiskPressure: shared.Fallback(diskPressure), Labels: shared.CloneStringMap(node.Labels), Conditions: conditions, Timeline: timeline}, nil
}

func List(clientset *kubernetes.Clientset) ([]ListItem, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]ListItem, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		items = append(items, ListItem{Name: node.Name, Status: status(node), Role: role(node), Version: shared.Fallback(node.Status.NodeInfo.KubeletVersion), InternalIP: address(node, corev1.NodeInternalIP), Age: shared.FormatAge(node.CreationTimestamp)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func YAML(clientset *kubernetes.Clientset, name string) (string, error) {
	node, err := clientset.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	content, err := yaml.Marshal(node)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ready(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func role(node corev1.Node) string {
	for key := range node.Labels {
		if strings.HasPrefix(key, "node-role.kubernetes.io/") {
			resourceRole := strings.TrimPrefix(key, "node-role.kubernetes.io/")
			if resourceRole == "" {
				return "Control plane"
			}
			return strings.ToUpper(resourceRole[:1]) + resourceRole[1:]
		}
	}
	return "Worker"
}

func status(node corev1.Node) string {
	if ready(node) {
		return "Ready"
	}
	if node.Spec.Unschedulable {
		return "Cordoned"
	}
	return "Not ready"
}

func address(node corev1.Node, addressType corev1.NodeAddressType) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == addressType {
			return shared.Fallback(addr.Address)
		}
	}
	return "-"
}

func boolStateLabel(ok bool) string {
	if ok {
		return "Clear"
	}
	return "Detected"
}
