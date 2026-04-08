package cluster

import (
	"context"
	"fmt"
	"net/http"

	"netkube/adapters"
	nodeapi "netkube/adapters/api/cluster/v1/node"
	"netkube/adapters/api/shared"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Handler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}
	nodes, err := nodeapi.NodesSummary(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	pvMetric, err := persistentVolumeMetric(cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	customResourceMetric, err := customResourceMetric(cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	usage, err := resourceUsage(cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Meta: shared.PageMetaFromCluster(cluster, ""), Nodes: Metric{Total: nodes.Total, Primary: nodes.Ready, Warning: nodes.Pending, Danger: nodes.NotReady, Status: nodes.Status}, PersistentVolumes: pvMetric, CustomResources: customResourceMetric, ResourceUsage: usage, Warnings: shared.ListWarningEvents(cluster.Clientset, "", 5)})
}

func persistentVolumeMetric(cluster *adapters.ClusterClient) (Metric, error) {
	volumes, err := cluster.Clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return Metric{}, err
	}

	metric := Metric{Total: len(volumes.Items)}
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

func customResourceMetric(cluster *adapters.ClusterClient) (Metric, error) {
	crds, err := cluster.APIExtensionsClient.ApiextensionsV1().CustomResourceDefinitions().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return Metric{}, err
	}

	metric := Metric{Total: len(crds.Items)}
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

func resourceUsage(cluster *adapters.ClusterClient) (ResourceUsage, error) {
	nodes, err := cluster.Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return ResourceUsage{}, err
	}

	pods, err := cluster.Clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return ResourceUsage{}, err
	}

	var usageCPUMilli, usageMemoryBytes int64
	nodeMetrics, err := cluster.MetricsClient.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err == nil {
		for _, item := range nodeMetrics.Items {
			usageCPUMilli += item.Usage.Cpu().MilliValue()
			usageMemoryBytes += item.Usage.Memory().Value()
		}
	}

	var capacityCPUMilli, capacityMemoryBytes, allocatablePods, allocatableCPUMilli, allocatableMemoryBytes int64
	for _, node := range nodes.Items {
		capacityCPUMilli += node.Status.Capacity.Cpu().MilliValue()
		capacityMemoryBytes += node.Status.Capacity.Memory().Value()
		allocatablePods += node.Status.Allocatable.Pods().Value()
		allocatableCPUMilli += node.Status.Allocatable.Cpu().MilliValue()
		allocatableMemoryBytes += node.Status.Allocatable.Memory().Value()
	}

	var requestedPods, requestedCPUMilli, requestedMemoryBytes int64
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		requestedPods++
		requestedCPUMilli += podRequestedCPUMilli(pod)
		requestedMemoryBytes += podRequestedMemoryBytes(pod)
	}

	return ResourceUsage{
		UsageCapacity:    ResourceUsageSection{CPU: ResourceUsageMetric{Percent: shared.PercentageFloat(usageCPUMilli, capacityCPUMilli), Used: shared.FormatCPU(usageCPUMilli), Total: shared.FormatCPU(capacityCPUMilli)}, Memory: ResourceUsageMetric{Percent: shared.PercentageFloat(usageMemoryBytes, capacityMemoryBytes), Used: shared.FormatBytes(usageMemoryBytes), Total: shared.FormatBytes(capacityMemoryBytes)}},
		RequestsAllocate: ResourceUsageSection{Pods: ResourceUsageMetric{Percent: shared.PercentageFloat(requestedPods, allocatablePods), Used: fmt.Sprintf("%d", requestedPods), Total: fmt.Sprintf("%d", allocatablePods)}, CPU: ResourceUsageMetric{Percent: shared.PercentageFloat(requestedCPUMilli, allocatableCPUMilli), Used: shared.FormatCPU(requestedCPUMilli), Total: shared.FormatCPU(allocatableCPUMilli)}, Memory: ResourceUsageMetric{Percent: shared.PercentageFloat(requestedMemoryBytes, allocatableMemoryBytes), Used: shared.FormatBytes(requestedMemoryBytes), Total: shared.FormatBytes(allocatableMemoryBytes)}},
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
	var initMax, total int64
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
