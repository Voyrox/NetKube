package api

import (
	"context"
	"fmt"
	"net/http"

	"netkube/adapters"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type overviewMetric struct {
	Total   int    `json:"total"`
	Primary int    `json:"primary"`
	Warning int    `json:"warning"`
	Danger  int    `json:"danger"`
	Status  string `json:"status"`
}

type clusterOverviewResponse struct {
	Meta              pageMeta       `json:"meta"`
	Nodes             overviewMetric `json:"nodes"`
	PersistentVolumes overviewMetric `json:"persistentVolumes"`
	CustomResources   overviewMetric `json:"customResources"`
	ResourceUsage     resourceUsage  `json:"resourceUsage"`
	Warnings          []warningEvent `json:"warnings"`
}

type resourceUsageMetric struct {
	Percent float64 `json:"percent"`
	Used    string  `json:"used"`
	Total   string  `json:"total"`
}

type resourceUsageSection struct {
	CPU    resourceUsageMetric `json:"cpu"`
	Memory resourceUsageMetric `json:"memory"`
	Pods   resourceUsageMetric `json:"pods,omitempty"`
}

type resourceUsage struct {
	UsageCapacity    resourceUsageSection `json:"usageCapacity"`
	RequestsAllocate resourceUsageSection `json:"requestsAllocate"`
}

type workloadsOverviewResponse struct {
	Meta           pageMeta       `json:"meta"`
	Pods           overviewMetric `json:"pods"`
	Deployments    overviewMetric `json:"deployments"`
	ReplicaSets    overviewMetric `json:"replicaSets"`
	DaemonSets     overviewMetric `json:"daemonSets"`
	StatefulSets   overviewMetric `json:"statefulSets"`
	CronJobs       overviewMetric `json:"cronJobs"`
	Jobs           overviewMetric `json:"jobs"`
	ResourceQuotas overviewMetric `json:"resourceQuotas"`
	Warnings       []warningEvent `json:"warnings"`
}

func ClusterOverviewHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	nodes, err := GetNodes(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	pvMetric, err := getPersistentVolumeMetric(cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	crdMetric, err := getCRDMetric(cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	resourceUsage, err := getClusterResourceUsage(cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, clusterOverviewResponse{
		Meta: pageMetaFromCluster(cluster, ""),
		Nodes: overviewMetric{
			Total:   nodes.Total,
			Primary: nodes.Ready,
			Warning: nodes.Pending,
			Danger:  nodes.NotReady,
			Status:  nodes.Status,
		},
		PersistentVolumes: pvMetric,
		CustomResources:   crdMetric,
		ResourceUsage:     resourceUsage,
		Warnings:          listWarningEvents(cluster.Clientset, "", 5),
	})
}

func WorkloadsOverviewHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	pods, podStats, err := GetPods(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	deployments, deploymentStats, err := GetDeployments(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	replicaSets, err := cluster.Clientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	daemonSets, err := cluster.Clientset.AppsV1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	statefulSets, err := cluster.Clientset.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	jobs, err := cluster.Clientset.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	cronJobs, err := cluster.Clientset.BatchV1().CronJobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	resourceQuotas, err := cluster.Clientset.CoreV1().ResourceQuotas(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, workloadsOverviewResponse{
		Meta: pageMetaFromCluster(cluster, namespace),
		Pods: overviewMetric{
			Total:   len(pods),
			Primary: podStats.Running,
			Warning: podStats.Pending,
			Danger:  podStats.Failed,
			Status:  statusFromCounts(podStats.Failed, podStats.Pending),
		},
		Deployments: overviewMetric{
			Total:   len(deployments),
			Primary: deploymentStats.Healthy,
			Warning: deploymentStats.Warning,
			Danger:  deploymentStats.Pending,
			Status:  statusFromCounts(deploymentStats.Pending, deploymentStats.Warning),
		},
		ReplicaSets:    simpleOverviewMetric(len(replicaSets.Items), readyReplicaSets(replicaSets.Items), pendingReplicaSets(replicaSets.Items), len(replicaSets.Items)-readyReplicaSets(replicaSets.Items)-pendingReplicaSets(replicaSets.Items)),
		DaemonSets:     simpleOverviewMetric(len(daemonSets.Items), readyDaemonSets(daemonSets.Items), pendingDaemonSets(daemonSets.Items), len(daemonSets.Items)-readyDaemonSets(daemonSets.Items)-pendingDaemonSets(daemonSets.Items)),
		StatefulSets:   simpleOverviewMetric(len(statefulSets.Items), readyStatefulSets(statefulSets.Items), updatingStatefulSets(statefulSets.Items), len(statefulSets.Items)-readyStatefulSets(statefulSets.Items)-updatingStatefulSets(statefulSets.Items)),
		CronJobs:       simpleOverviewMetric(len(cronJobs.Items), len(cronJobs.Items), 0, 0),
		Jobs:           simpleOverviewMetric(len(jobs.Items), successfulJobs(jobs.Items), activeJobs(jobs.Items), failedJobs(jobs.Items)),
		ResourceQuotas: simpleOverviewMetric(len(resourceQuotas.Items), len(resourceQuotas.Items), 0, 0),
		Warnings:       listWarningEvents(cluster.Clientset, namespace, 5),
	})
}

func getPersistentVolumeMetric(cluster *adapters.ClusterClient) (overviewMetric, error) {
	volumes, err := cluster.Clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return overviewMetric{}, err
	}

	metric := overviewMetric{Total: len(volumes.Items)}
	for _, item := range volumes.Items {
		switch item.Status.Phase {
		case corev1.VolumeBound:
			metric.Primary++
		case corev1.VolumePending:
			metric.Warning++
		case corev1.VolumeFailed:
			metric.Danger++
		}
	}

	metric.Status = statusFromCounts(metric.Danger, metric.Warning)
	return metric, nil
}

func getCRDMetric(cluster *adapters.ClusterClient) (overviewMetric, error) {
	crds, err := cluster.APIExtensionsClient.ApiextensionsV1().CustomResourceDefinitions().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return overviewMetric{}, err
	}

	metric := overviewMetric{Total: len(crds.Items)}
	for _, item := range crds.Items {
		established := false
		terminating := false
		for _, condition := range item.Status.Conditions {
			if condition.Type == "Established" && condition.Status == "True" {
				established = true
			}
			if condition.Type == "Terminating" && condition.Status == "True" {
				terminating = true
			}
		}

		switch {
		case terminating:
			metric.Danger++
		case established:
			metric.Primary++
		default:
			metric.Warning++
		}
	}

	metric.Status = statusFromCounts(metric.Danger, metric.Warning)
	return metric, nil
}

func statusFromCounts(danger, warning int) string {
	switch {
	case danger > 0:
		return "Attention"
	case warning > 0:
		return "Watch"
	default:
		return "Healthy"
	}
}

func simpleOverviewMetric(total, primary, warning, danger int) overviewMetric {
	return overviewMetric{
		Total:   total,
		Primary: primary,
		Warning: warning,
		Danger:  danger,
		Status:  statusFromCounts(danger, warning),
	}
}

func getClusterResourceUsage(cluster *adapters.ClusterClient) (resourceUsage, error) {
	nodes, err := cluster.Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return resourceUsage{}, err
	}

	pods, err := cluster.Clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return resourceUsage{}, err
	}

	var usageCPUMilli int64
	var usageMemoryBytes int64
	nodeMetrics, err := cluster.MetricsClient.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err == nil {
		for _, item := range nodeMetrics.Items {
			usageCPUMilli += item.Usage.Cpu().MilliValue()
			usageMemoryBytes += item.Usage.Memory().Value()
		}
	}

	var capacityCPUMilli int64
	var capacityMemoryBytes int64
	var allocatablePods int64
	var allocatableCPUMilli int64
	var allocatableMemoryBytes int64

	for _, node := range nodes.Items {
		capacityCPUMilli += node.Status.Capacity.Cpu().MilliValue()
		capacityMemoryBytes += node.Status.Capacity.Memory().Value()
		allocatablePods += node.Status.Allocatable.Pods().Value()
		allocatableCPUMilli += node.Status.Allocatable.Cpu().MilliValue()
		allocatableMemoryBytes += node.Status.Allocatable.Memory().Value()
	}

	var requestedPods int64
	var requestedCPUMilli int64
	var requestedMemoryBytes int64
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		requestedPods++
		requestedCPUMilli += podRequestedCPUMilli(pod)
		requestedMemoryBytes += podRequestedMemoryBytes(pod)
	}

	return resourceUsage{
		UsageCapacity: resourceUsageSection{
			CPU: resourceUsageMetric{
				Percent: percentageFloat(usageCPUMilli, capacityCPUMilli),
				Used:    formatCPU(usageCPUMilli),
				Total:   formatCPU(capacityCPUMilli),
			},
			Memory: resourceUsageMetric{
				Percent: percentageFloat(usageMemoryBytes, capacityMemoryBytes),
				Used:    formatBytes(usageMemoryBytes),
				Total:   formatBytes(capacityMemoryBytes),
			},
		},
		RequestsAllocate: resourceUsageSection{
			Pods: resourceUsageMetric{
				Percent: percentageFloat(requestedPods, allocatablePods),
				Used:    fmt.Sprintf("%d", requestedPods),
				Total:   fmt.Sprintf("%d", allocatablePods),
			},
			CPU: resourceUsageMetric{
				Percent: percentageFloat(requestedCPUMilli, allocatableCPUMilli),
				Used:    formatCPU(requestedCPUMilli),
				Total:   formatCPU(allocatableCPUMilli),
			},
			Memory: resourceUsageMetric{
				Percent: percentageFloat(requestedMemoryBytes, allocatableMemoryBytes),
				Used:    formatBytes(requestedMemoryBytes),
				Total:   formatBytes(allocatableMemoryBytes),
			},
		},
	}, nil
}

func podRequestedCPUMilli(pod corev1.Pod) int64 {
	var total int64
	for _, container := range pod.Spec.InitContainers {
		if value := container.Resources.Requests.Cpu(); value != nil && value.MilliValue() > total {
			total = value.MilliValue()
		}
	}
	for _, container := range pod.Spec.Containers {
		if value := container.Resources.Requests.Cpu(); value != nil {
			total += value.MilliValue()
		}
	}
	if pod.Spec.Overhead != nil {
		if value := pod.Spec.Overhead.Cpu(); value != nil {
			total += value.MilliValue()
		}
	}
	return total
}

func podRequestedMemoryBytes(pod corev1.Pod) int64 {
	var initMax int64
	var total int64
	for _, container := range pod.Spec.InitContainers {
		if value := container.Resources.Requests.Memory(); value != nil && value.Value() > initMax {
			initMax = value.Value()
		}
	}
	for _, container := range pod.Spec.Containers {
		if value := container.Resources.Requests.Memory(); value != nil {
			total += value.Value()
		}
	}
	if pod.Spec.Overhead != nil {
		if value := pod.Spec.Overhead.Memory(); value != nil {
			total += value.Value()
		}
	}
	if initMax > total {
		return initMax
	}
	return total
}

func percentageFloat(used, total int64) float64 {
	if total <= 0 {
		return 0
	}
	return (float64(used) / float64(total)) * 100
}

func formatCPU(milli int64) string {
	if milli < 1000 {
		return fmt.Sprintf("%dm", milli)
	}
	value := float64(milli) / 1000
	return trimFloat(value) + " cores"
}

func formatBytes(bytes int64) string {
	const (
		ki = 1024
		mi = ki * 1024
		gi = mi * 1024
	)

	switch {
	case bytes >= gi:
		return trimFloat(float64(bytes)/float64(gi)) + "Gi"
	case bytes >= mi:
		return trimFloat(float64(bytes)/float64(mi)) + "Mi"
	case bytes >= ki:
		return trimFloat(float64(bytes)/float64(ki)) + "Ki"
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

func trimFloat(value float64) string {
	formatted := fmt.Sprintf("%.2f", value)
	for len(formatted) > 0 && formatted[len(formatted)-1] == '0' {
		formatted = formatted[:len(formatted)-1]
	}
	if len(formatted) > 0 && formatted[len(formatted)-1] == '.' {
		formatted = formatted[:len(formatted)-1]
	}
	return formatted
}

func readyReplicaSets(items []appsv1.ReplicaSet) int {
	count := 0
	for _, item := range items {
		if item.Status.ReadyReplicas == desiredReplicas(item.Spec.Replicas) {
			count++
		}
	}
	return count
}

func pendingReplicaSets(items []appsv1.ReplicaSet) int {
	count := 0
	for _, item := range items {
		if item.Status.Replicas == 0 || item.Status.ReadyReplicas == 0 {
			count++
		}
	}
	return count
}

func readyDaemonSets(items []appsv1.DaemonSet) int {
	count := 0
	for _, item := range items {
		if item.Status.DesiredNumberScheduled > 0 && item.Status.NumberReady == item.Status.DesiredNumberScheduled {
			count++
		}
	}
	return count
}

func pendingDaemonSets(items []appsv1.DaemonSet) int {
	count := 0
	for _, item := range items {
		if item.Status.DesiredNumberScheduled == 0 || item.Status.NumberReady == 0 {
			count++
		}
	}
	return count
}

func readyStatefulSets(items []appsv1.StatefulSet) int {
	count := 0
	for _, item := range items {
		if item.Status.AvailableReplicas == desiredReplicas(item.Spec.Replicas) {
			count++
		}
	}
	return count
}

func updatingStatefulSets(items []appsv1.StatefulSet) int {
	count := 0
	for _, item := range items {
		if item.Status.UpdatedReplicas < item.Status.Replicas {
			count++
		}
	}
	return count
}

func successfulJobs(items []batchv1.Job) int {
	count := 0
	for _, item := range items {
		if item.Status.Succeeded > 0 {
			count++
		}
	}
	return count
}

func activeJobs(items []batchv1.Job) int {
	count := 0
	for _, item := range items {
		if item.Status.Active > 0 {
			count++
		}
	}
	return count
}

func failedJobs(items []batchv1.Job) int {
	count := 0
	for _, item := range items {
		if item.Status.Failed > 0 {
			count++
		}
	}
	return count
}

func desiredReplicas(value *int32) int32 {
	if value == nil {
		return 1
	}

	return *value
}
