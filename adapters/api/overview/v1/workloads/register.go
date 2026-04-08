package workloads

import (
	"context"
	"net/http"

	"netkube/adapters/api/shared"
	cronjobapi "netkube/adapters/api/workloads/v1/cronjob"
	daemonsetapi "netkube/adapters/api/workloads/v1/daemonset"
	deploymentapi "netkube/adapters/api/workloads/v1/deployment"
	jobapi "netkube/adapters/api/workloads/v1/job"
	podapi "netkube/adapters/api/workloads/v1/pod"
	replicasetapi "netkube/adapters/api/workloads/v1/replicaset"
	statefulsetapi "netkube/adapters/api/workloads/v1/statefulset"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Handler(c *gin.Context) {
	cluster, ok := shared.ResolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := shared.SelectedNamespace(c)
	pods, podStats, err := podapi.List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	deployments, deploymentStats, err := deploymentapi.List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	replicaSets, replicaSetStats, err := replicasetapi.List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	daemonSets, daemonSetStats, err := daemonsetapi.List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	statefulSets, statefulSetStats, err := statefulsetapi.List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	jobs, jobStats, err := jobapi.List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	cronJobs, cronJobStats, err := cronjobapi.List(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	resourceQuotas, err := cluster.Clientset.CoreV1().ResourceQuotas(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, shared.APIError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{
		Meta:           shared.PageMetaFromCluster(cluster, namespace),
		Pods:           Metric{Total: len(pods), Primary: podStats.Running, Warning: podStats.Pending, Danger: podStats.Failed, Status: statusFromCounts(podStats.Failed, podStats.Pending)},
		Deployments:    Metric{Total: len(deployments), Primary: deploymentStats.Healthy, Warning: deploymentStats.Warning, Danger: deploymentStats.Pending, Status: statusFromCounts(deploymentStats.Pending, deploymentStats.Warning)},
		ReplicaSets:    Metric{Total: len(replicaSets), Primary: replicaSetStats.Healthy, Warning: replicaSetStats.Warning, Danger: replicaSetStats.Pending, Status: statusFromCounts(replicaSetStats.Pending, replicaSetStats.Warning)},
		DaemonSets:     Metric{Total: len(daemonSets), Primary: daemonSetStats.Healthy, Warning: daemonSetStats.Warning, Danger: daemonSetStats.Pending, Status: statusFromCounts(daemonSetStats.Pending, daemonSetStats.Warning)},
		StatefulSets:   Metric{Total: len(statefulSets), Primary: statefulSetStats.Healthy, Warning: statefulSetStats.Warning, Danger: statefulSetStats.Pending, Status: statusFromCounts(statefulSetStats.Pending, statefulSetStats.Warning)},
		CronJobs:       Metric{Total: len(cronJobs), Primary: cronJobStats.Scheduled, Warning: cronJobStats.Suspended, Danger: cronJobStats.Active, Status: statusFromCounts(cronJobStats.Active, cronJobStats.Suspended)},
		Jobs:           Metric{Total: len(jobs), Primary: jobStats.Succeeded, Warning: jobStats.Active, Danger: jobStats.Failed, Status: statusFromCounts(jobStats.Failed, jobStats.Active)},
		ResourceQuotas: Metric{Total: len(resourceQuotas.Items), Primary: len(resourceQuotas.Items), Status: statusFromCounts(0, 0)},
		Warnings:       shared.ListWarningEvents(cluster.Clientset, namespace, 5),
	})

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
