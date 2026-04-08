package pod

import (
	"context"
	"io"
	"sort"
	"strings"

	"netkube/adapters/api/shared"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func List(clientset *kubernetes.Clientset, namespace string) ([]Row, Stats, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, Stats{}, err
	}

	items := make([]Row, 0, len(pods.Items))
	stats := Stats{Total: len(pods.Items)}

	for _, item := range pods.Items {
		status := status(item)
		switch status {
		case "Running":
			stats.Running++
		case "Pending":
			stats.Pending++
		case "Failed":
			stats.Failed++
		default:
			stats.Other++
		}

		restarts, lastRestart, lastRestartReason := restartInfo(item)
		items = append(items, Row{
			Namespace:         item.Namespace,
			Name:              item.Name,
			Ready:             shared.ReadyRatio(readyContainers(item), int32(len(item.Status.ContainerStatuses))),
			Status:            status,
			Restarts:          restarts,
			LastRestart:       lastRestart,
			LastRestartReason: lastRestartReason,
			Node:              shared.Fallback(item.Spec.NodeName),
			PodIP:             shared.Fallback(item.Status.PodIP),
			Age:               shared.FormatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		}

		return strings.ToLower(items[i].Namespace) < strings.ToLower(items[j].Namespace)
	})

	return items, stats, nil
}

func Create(clientset *kubernetes.Clientset, content string) (*corev1.Pod, error) {
	pod, err := decodeManifest(content)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Pods(pod.Namespace).Create(context.Background(), pod, metav1.CreateOptions{})
}

func DetailFor(clientset *kubernetes.Clientset, namespace, name string) (Detail, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return Detail{}, err
	}

	restarts, lastRestart, lastRestartReason := restartInfo(*pod)
	containers := make([]ContainerRow, 0, len(pod.Status.ContainerStatuses))
	for _, containerStatus := range pod.Status.ContainerStatuses {
		containers = append(containers, ContainerRow{
			Name:     containerStatus.Name,
			Image:    imageForContainer(*pod, containerStatus.Name),
			Ready:    containerStatus.Ready,
			Restarts: containerStatus.RestartCount,
			State:    containerState(containerStatus),
		})
	}

	conditions := make([]ConditionRow, 0, len(pod.Status.Conditions))
	for _, condition := range pod.Status.Conditions {
		conditions = append(conditions, ConditionRow{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  shared.Fallback(condition.Reason),
			Message: shared.Fallback(condition.Message),
		})
	}

	return Detail{
		Namespace:         pod.Namespace,
		Name:              pod.Name,
		Ready:             shared.ReadyRatio(readyContainers(*pod), int32(len(pod.Status.ContainerStatuses))),
		Status:            status(*pod),
		Phase:             shared.Fallback(string(pod.Status.Phase)),
		Restarts:          restarts,
		LastRestart:       lastRestart,
		LastRestartReason: lastRestartReason,
		Node:              shared.Fallback(pod.Spec.NodeName),
		PodIP:             shared.Fallback(pod.Status.PodIP),
		HostIP:            shared.Fallback(pod.Status.HostIP),
		ServiceAccount:    shared.Fallback(pod.Spec.ServiceAccountName),
		QOSClass:          shared.Fallback(string(pod.Status.QOSClass)),
		Age:               shared.FormatAge(pod.CreationTimestamp),
		StartTime:         shared.FormatTime(pod.Status.StartTime),
		Labels:            shared.CloneStringMap(pod.Labels),
		Annotations:       shared.CloneStringMap(pod.Annotations),
		Containers:        containers,
		Conditions:        conditions,
	}, nil
}

func Logs(clientset *kubernetes.Clientset, namespace, name, container string) (string, string, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	selectedContainer := container
	if selectedContainer == "" {
		if len(pod.Spec.Containers) == 0 {
			return "", "", nil
		}
		selectedContainer = pod.Spec.Containers[0].Name
	}

	tailLines := int64(200)
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{Container: selectedContainer, TailLines: &tailLines})
	stream, err := req.Stream(context.Background())
	if err != nil {
		return "", selectedContainer, err
	}
	defer stream.Close()

	content, err := io.ReadAll(stream)
	if err != nil {
		return "", selectedContainer, err
	}

	return string(content), selectedContainer, nil
}

func Events(clientset *kubernetes.Clientset, namespace, name string) ([]EventRow, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{FieldSelector: "involvedObject.kind=Pod,involvedObject.name=" + name})
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
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	content, err := yaml.Marshal(pod)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func status(pod corev1.Pod) string {
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

func containerState(status corev1.ContainerStatus) string {
	switch {
	case status.State.Running != nil:
		return "Running"
	case status.State.Waiting != nil:
		return shared.Fallback(status.State.Waiting.Reason)
	case status.State.Terminated != nil:
		return shared.Fallback(status.State.Terminated.Reason)
	default:
		return "Unknown"
	}
}

func imageForContainer(pod corev1.Pod, name string) string {
	for _, container := range pod.Spec.InitContainers {
		if container.Name == name {
			return shared.Fallback(container.Image)
		}
	}
	for _, container := range pod.Spec.Containers {
		if container.Name == name {
			return shared.Fallback(container.Image)
		}
	}
	return "-"
}

func readyContainers(pod corev1.Pod) int32 {
	var ready int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			ready++
		}
	}
	return ready
}

func restartInfo(pod corev1.Pod) (int32, string, string) {
	var restarts int32
	lastRestart := "-"
	lastRestartReason := "-"
	var latest metav1.Time

	for _, containerStatus := range pod.Status.ContainerStatuses {
		restarts += containerStatus.RestartCount
		if containerStatus.LastTerminationState.Terminated == nil {
			continue
		}

		terminated := containerStatus.LastTerminationState.Terminated
		if latest.IsZero() || terminated.FinishedAt.After(latest.Time) {
			latest = terminated.FinishedAt
			lastRestart = shared.FormatAge(terminated.FinishedAt)
			lastRestartReason = shared.Fallback(terminated.Reason)
		}
	}

	return restarts, lastRestart, lastRestartReason
}
