package api

import (
	"context"
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
	Warnings          []warningEvent `json:"warnings"`
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
