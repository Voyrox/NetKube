package deployment

import (
	"context"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func List(clientset *kubernetes.Clientset, namespace string) ([]Row, Stats, error) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, Stats{}, err
	}

	items := make([]Row, 0, len(deployments.Items))
	stats := Stats{Total: len(deployments.Items)}

	for _, item := range deployments.Items {
		status := status(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		items = append(items, Row{Namespace: item.Namespace, Name: item.Name, Ready: shared.ReadyRatio(item.Status.ReadyReplicas, item.Status.Replicas), Status: status, Desired: item.Status.Replicas, Updated: item.Status.UpdatedReplicas, Available: item.Status.AvailableReplicas, Age: shared.FormatAge(item.CreationTimestamp)})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		}
		return strings.ToLower(items[i].Namespace) < strings.ToLower(items[j].Namespace)
	})

	return items, stats, nil
}

func Create(clientset *kubernetes.Clientset, content string) (*appsv1.Deployment, error) {
	deployment, err := decodeManifest(content)
	if err != nil {
		return nil, err
	}

	return clientset.AppsV1().Deployments(deployment.Namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
}

func DetailFor(clientset *kubernetes.Clientset, namespace, name string) (Detail, error) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return Detail{}, err
	}

	selector := labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels)
	selectorString := selector.String()

	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: selectorString})
	if err != nil {
		return Detail{}, err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: selectorString})
	if err != nil {
		return Detail{}, err
	}

	conditions := make([]ConditionRow, 0, len(deployment.Status.Conditions))
	for _, condition := range deployment.Status.Conditions {
		conditions = append(conditions, ConditionRow{Type: string(condition.Type), Status: string(condition.Status), Reason: shared.Fallback(condition.Reason), Message: shared.Fallback(condition.Message)})
	}

	replicaSetItems := make([]ReplicaSetRow, 0, len(replicaSets.Items))
	for _, item := range replicaSets.Items {
		replicaSetItems = append(replicaSetItems, ReplicaSetRow{Name: item.Name, Ready: shared.ReadyRatio(item.Status.ReadyReplicas, item.Status.Replicas), Desired: item.Status.Replicas, Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(replicaSetItems, func(i, j int) bool { return replicaSetItems[i].Name < replicaSetItems[j].Name })

	podItems := make([]PodRow, 0, len(pods.Items))
	for _, item := range pods.Items {
		podItems = append(podItems, PodRow{Name: item.Name, Ready: shared.ReadyRatio(podReadyContainers(item), int32(len(item.Status.ContainerStatuses))), Status: podStatus(item), Node: shared.Fallback(item.Spec.NodeName), Age: shared.FormatAge(item.CreationTimestamp)})
	}
	sort.Slice(podItems, func(i, j int) bool { return podItems[i].Name < podItems[j].Name })

	return Detail{
		Namespace:   deployment.Namespace,
		Name:        deployment.Name,
		Status:      status(*deployment),
		Ready:       shared.ReadyRatio(deployment.Status.ReadyReplicas, deployment.Status.Replicas),
		Desired:     deployment.Status.Replicas,
		Updated:     deployment.Status.UpdatedReplicas,
		Available:   deployment.Status.AvailableReplicas,
		Unavailable: deployment.Status.UnavailableReplicas,
		Age:         shared.FormatAge(deployment.CreationTimestamp),
		Strategy:    string(deployment.Spec.Strategy.Type),
		Selector:    shared.Fallback(selectorString),
		Conditions:  conditions,
		ReplicaSets: replicaSetItems,
		Pods:        podItems,
		Labels:      shared.CloneStringMap(deployment.Labels),
		Annotations: shared.CloneStringMap(deployment.Annotations),
	}, nil
}

func Events(clientset *kubernetes.Clientset, namespace, name string) ([]EventRow, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{FieldSelector: "involvedObject.kind=Deployment,involvedObject.name=" + name})
	if err != nil {
		return nil, err
	}

	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].CreationTimestamp.Time.After(events.Items[j].CreationTimestamp.Time)
	})

	items := make([]EventRow, 0, len(events.Items))
	for _, item := range events.Items {
		items = append(items, EventRow{Type: shared.Fallback(item.Type), Reason: shared.Fallback(item.Reason), Message: shared.Fallback(item.Message), Age: shared.FormatAge(item.CreationTimestamp)})
	}

	return items, nil
}

func YAML(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	content, err := yaml.Marshal(deployment)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func status(item appsv1.Deployment) string {
	switch {
	case item.Status.Replicas == 0:
		return "Pending"
	case item.Status.AvailableReplicas == item.Status.Replicas && item.Status.UpdatedReplicas == item.Status.Replicas:
		return "Healthy"
	case item.Status.UnavailableReplicas > 0:
		return "Degraded"
	default:
		return "Updating"
	}
}

func podStatus(pod corev1.Pod) string {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil && containerStatus.State.Waiting.Reason != "" {
			return containerStatus.State.Waiting.Reason
		}
		if containerStatus.State.Terminated != nil && containerStatus.State.Terminated.Reason != "" {
			return containerStatus.State.Terminated.Reason
		}
	}

	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}
	if pod.Status.Phase == "" {
		return "Unknown"
	}

	return string(pod.Status.Phase)
}

func podReadyContainers(pod corev1.Pod) int32 {
	var ready int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			ready++
		}
	}
	return ready
}
