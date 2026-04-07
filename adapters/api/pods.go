package api

import (
	"context"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

type podRow struct {
	Namespace         string `json:"namespace"`
	Name              string `json:"name"`
	Ready             string `json:"ready"`
	Status            string `json:"status"`
	Restarts          int32  `json:"restarts"`
	LastRestart       string `json:"lastRestart"`
	LastRestartReason string `json:"lastRestartReason"`
	Node              string `json:"node"`
	PodIP             string `json:"podIP"`
	Age               string `json:"age"`
}

type podsStats struct {
	Running int `json:"running"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
	Other   int `json:"other"`
	Total   int `json:"total"`
}

type podsResponse struct {
	Meta  pageMeta  `json:"meta"`
	Items []podRow  `json:"items"`
	Count int       `json:"count"`
	Stats podsStats `json:"stats"`
}

type podContainerRow struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Ready    bool   `json:"ready"`
	Restarts int32  `json:"restarts"`
	State    string `json:"state"`
}

type podConditionRow struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type podDetail struct {
	Namespace         string            `json:"namespace"`
	Name              string            `json:"name"`
	Ready             string            `json:"ready"`
	Status            string            `json:"status"`
	Phase             string            `json:"phase"`
	Restarts          int32             `json:"restarts"`
	LastRestart       string            `json:"lastRestart"`
	LastRestartReason string            `json:"lastRestartReason"`
	Node              string            `json:"node"`
	PodIP             string            `json:"podIP"`
	HostIP            string            `json:"hostIP"`
	ServiceAccount    string            `json:"serviceAccount"`
	QOSClass          string            `json:"qosClass"`
	Age               string            `json:"age"`
	StartTime         string            `json:"startTime"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	Containers        []podContainerRow `json:"containers"`
	Conditions        []podConditionRow `json:"conditions"`
}

type podDetailResponse struct {
	Meta pageMeta  `json:"meta"`
	Item podDetail `json:"item"`
}

type podLogResponse struct {
	Meta      pageMeta `json:"meta"`
	Container string   `json:"container"`
	Content   string   `json:"content"`
}

type podEventRow struct {
	Type    string `json:"type"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Age     string `json:"age"`
}

type podEventsResponse struct {
	Meta  pageMeta      `json:"meta"`
	Items []podEventRow `json:"items"`
}

type podYAMLResponse struct {
	Meta    pageMeta `json:"meta"`
	Content string   `json:"content"`
}

func PodsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetPods(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, podsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func PodDetailHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "pod name and namespace are required"})
		return
	}

	item, err := GetPodDetail(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, podDetailResponse{
		Meta: pageMetaFromCluster(cluster, namespace),
		Item: item,
	})
}

func PodLogsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	container := strings.TrimSpace(c.Query("container"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "pod name and namespace are required"})
		return
	}

	content, selectedContainer, err := GetPodLogs(cluster.Clientset, namespace, name, container)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, podLogResponse{
		Meta:      pageMetaFromCluster(cluster, namespace),
		Container: selectedContainer,
		Content:   content,
	})
}

func PodEventsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "pod name and namespace are required"})
		return
	}

	items, err := GetPodEvents(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, podEventsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
	})
}

func PodYAMLHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "pod name and namespace are required"})
		return
	}

	content, err := GetPodYAML(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, podYAMLResponse{
		Meta:    pageMetaFromCluster(cluster, namespace),
		Content: content,
	})
}

func GetPods(clientset *kubernetes.Clientset, namespace string) ([]podRow, podsStats, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, podsStats{}, err
	}

	items := make([]podRow, 0, len(pods.Items))
	stats := podsStats{Total: len(pods.Items)}

	for _, item := range pods.Items {
		status := podStatus(item)
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

		restarts, lastRestart, lastRestartReason := containerRestartInfo(item)
		items = append(items, podRow{
			Namespace:         item.Namespace,
			Name:              item.Name,
			Ready:             readyRatio(readyContainers(item), int32(len(item.Status.ContainerStatuses))),
			Status:            status,
			Restarts:          restarts,
			LastRestart:       lastRestart,
			LastRestartReason: lastRestartReason,
			Node:              fallback(item.Spec.NodeName),
			PodIP:             fallback(item.Status.PodIP),
			Age:               formatAge(item.CreationTimestamp),
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

func GetPodDetail(clientset *kubernetes.Clientset, namespace, name string) (podDetail, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return podDetail{}, err
	}

	restarts, lastRestart, lastRestartReason := containerRestartInfo(*pod)
	containers := make([]podContainerRow, 0, len(pod.Status.ContainerStatuses))
	for _, status := range pod.Status.ContainerStatuses {
		containers = append(containers, podContainerRow{
			Name:     status.Name,
			Image:    containerImageForName(*pod, status.Name),
			Ready:    status.Ready,
			Restarts: status.RestartCount,
			State:    containerState(status),
		})
	}

	conditions := make([]podConditionRow, 0, len(pod.Status.Conditions))
	for _, condition := range pod.Status.Conditions {
		conditions = append(conditions, podConditionRow{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  fallback(condition.Reason),
			Message: fallback(condition.Message),
		})
	}

	startTime := "-"
	if pod.Status.StartTime != nil && !pod.Status.StartTime.IsZero() {
		startTime = pod.Status.StartTime.Time.Format("2006-01-02 15:04:05 MST")
	}

	return podDetail{
		Namespace:         pod.Namespace,
		Name:              pod.Name,
		Ready:             readyRatio(readyContainers(*pod), int32(len(pod.Status.ContainerStatuses))),
		Status:            podStatus(*pod),
		Phase:             fallback(string(pod.Status.Phase)),
		Restarts:          restarts,
		LastRestart:       lastRestart,
		LastRestartReason: lastRestartReason,
		Node:              fallback(pod.Spec.NodeName),
		PodIP:             fallback(pod.Status.PodIP),
		HostIP:            fallback(pod.Status.HostIP),
		ServiceAccount:    fallback(pod.Spec.ServiceAccountName),
		QOSClass:          fallback(string(pod.Status.QOSClass)),
		Age:               formatAge(pod.CreationTimestamp),
		StartTime:         startTime,
		Labels:            cloneStringMap(pod.Labels),
		Annotations:       cloneStringMap(pod.Annotations),
		Containers:        containers,
		Conditions:        conditions,
	}, nil
}

func GetPodLogs(clientset *kubernetes.Clientset, namespace, name, container string) (string, string, error) {
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
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
		Container: selectedContainer,
		TailLines: &tailLines,
	})

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

func GetPodEvents(clientset *kubernetes.Clientset, namespace, name string) ([]podEventRow, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{
		FieldSelector: "involvedObject.kind=Pod,involvedObject.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	items := make([]podEventRow, 0, len(events.Items))
	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].CreationTimestamp.Time.After(events.Items[j].CreationTimestamp.Time)
	})

	for _, item := range events.Items {
		items = append(items, podEventRow{
			Type:    fallback(item.Type),
			Reason:  fallback(item.Reason),
			Message: fallback(item.Message),
			Age:     formatAge(item.CreationTimestamp),
		})
	}

	return items, nil
}

func GetPodYAML(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
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

func podStatus(pod corev1.Pod) string {
	for _, status := range pod.Status.ContainerStatuses {
		if status.State.Waiting != nil && status.State.Waiting.Reason != "" {
			return status.State.Waiting.Reason
		}
		if status.State.Terminated != nil && status.State.Terminated.Reason != "" {
			return status.State.Terminated.Reason
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
		return fallback(status.State.Waiting.Reason)
	case status.State.Terminated != nil:
		return fallback(status.State.Terminated.Reason)
	default:
		return "Unknown"
	}
}

func containerImageForName(pod corev1.Pod, name string) string {
	for _, container := range pod.Spec.InitContainers {
		if container.Name == name {
			return fallback(container.Image)
		}
	}

	for _, container := range pod.Spec.Containers {
		if container.Name == name {
			return fallback(container.Image)
		}
	}

	return "-"
}

func cloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return map[string]string{}
	}

	clone := make(map[string]string, len(values))
	for key, value := range values {
		clone[key] = value
	}

	return clone
}

func readyContainers(pod corev1.Pod) int32 {
	var ready int32
	for _, status := range pod.Status.ContainerStatuses {
		if status.Ready {
			ready++
		}
	}

	return ready
}

func containerRestartInfo(pod corev1.Pod) (int32, string, string) {
	var restarts int32
	lastRestart := "-"
	lastRestartReason := "-"
	var latest metav1.Time

	for _, status := range pod.Status.ContainerStatuses {
		restarts += status.RestartCount
		if status.LastTerminationState.Terminated == nil {
			continue
		}

		terminated := status.LastTerminationState.Terminated
		if latest.IsZero() || terminated.FinishedAt.After(latest.Time) {
			latest = terminated.FinishedAt
			lastRestart = formatAge(terminated.FinishedAt)
			lastRestartReason = fallback(terminated.Reason)
		}
	}

	return restarts, lastRestart, lastRestartReason
}

func readyRatio(current, total int32) string {
	return strings.TrimSpace(strings.Join([]string{int32String(current), int32String(total)}, "/"))
}

func int32String(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func fallback(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}

	return value
}
