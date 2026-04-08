package api

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ReplicaSetsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetReplicaSets(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, replicaSetsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func DaemonSetsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetDaemonSets(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, daemonSetsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func StatefulSetsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetStatefulSets(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, statefulSetsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func JobsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetJobs(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func CronJobsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetCronJobs(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, cronJobsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func GetReplicaSets(clientset *kubernetes.Clientset, namespace string) ([]replicaSetRow, replicaSetsStats, error) {
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, replicaSetsStats{}, err
	}

	items := make([]replicaSetRow, 0, len(replicaSets.Items))
	stats := replicaSetsStats{Total: len(replicaSets.Items)}

	for _, item := range replicaSets.Items {
		status := replicaSetStatus(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		desired := replicasOrZero(item.Spec.Replicas)
		items = append(items, replicaSetRow{
			Namespace: item.Namespace,
			Name:      item.Name,
			Ready:     readyRatio(item.Status.ReadyReplicas, desired),
			Status:    status,
			Desired:   desired,
			Current:   item.Status.Replicas,
			ReadyPods: item.Status.ReadyReplicas,
			Age:       formatAge(item.CreationTimestamp),
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

func GetDaemonSets(clientset *kubernetes.Clientset, namespace string) ([]daemonSetRow, daemonSetsStats, error) {
	daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, daemonSetsStats{}, err
	}

	items := make([]daemonSetRow, 0, len(daemonSets.Items))
	stats := daemonSetsStats{Total: len(daemonSets.Items)}

	for _, item := range daemonSets.Items {
		status := daemonSetStatus(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		items = append(items, daemonSetRow{
			Namespace:    item.Namespace,
			Name:         item.Name,
			Ready:        readyRatio(item.Status.NumberReady, item.Status.DesiredNumberScheduled),
			Status:       status,
			Desired:      item.Status.DesiredNumberScheduled,
			Current:      item.Status.CurrentNumberScheduled,
			Available:    item.Status.NumberAvailable,
			Misscheduled: item.Status.NumberMisscheduled,
			Age:          formatAge(item.CreationTimestamp),
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

func GetStatefulSets(clientset *kubernetes.Clientset, namespace string) ([]statefulSetRow, statefulSetsStats, error) {
	statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, statefulSetsStats{}, err
	}

	items := make([]statefulSetRow, 0, len(statefulSets.Items))
	stats := statefulSetsStats{Total: len(statefulSets.Items)}

	for _, item := range statefulSets.Items {
		status := statefulSetStatus(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		desired := desiredReplicas(item.Spec.Replicas)
		items = append(items, statefulSetRow{
			Namespace: item.Namespace,
			Name:      item.Name,
			Ready:     readyRatio(item.Status.ReadyReplicas, desired),
			Status:    status,
			Desired:   desired,
			Current:   item.Status.CurrentReplicas,
			Updated:   item.Status.UpdatedReplicas,
			Age:       formatAge(item.CreationTimestamp),
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

func GetJobs(clientset *kubernetes.Clientset, namespace string) ([]jobRow, jobsStats, error) {
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, jobsStats{}, err
	}

	items := make([]jobRow, 0, len(jobs.Items))
	stats := jobsStats{Total: len(jobs.Items)}

	for _, item := range jobs.Items {
		status := jobStatus(item)
		switch status {
		case "Succeeded":
			stats.Succeeded++
		case "Running":
			stats.Active++
		case "Failed":
			stats.Failed++
		}

		items = append(items, jobRow{
			Namespace:   item.Namespace,
			Name:        item.Name,
			Status:      status,
			Completions: readyRatio(item.Status.Succeeded, desiredCompletions(item.Spec.Completions)),
			Active:      item.Status.Active,
			Duration:    jobDuration(item),
			Age:         formatAge(item.CreationTimestamp),
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

func GetCronJobs(clientset *kubernetes.Clientset, namespace string) ([]cronJobRow, cronJobsStats, error) {
	cronJobs, err := clientset.BatchV1().CronJobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, cronJobsStats{}, err
	}

	items := make([]cronJobRow, 0, len(cronJobs.Items))
	stats := cronJobsStats{Total: len(cronJobs.Items)}

	for _, item := range cronJobs.Items {
		status := cronJobStatus(item)
		if item.Spec.Suspend != nil && *item.Spec.Suspend {
			stats.Suspended++
		} else {
			stats.Scheduled++
		}
		if len(item.Status.Active) > 0 {
			stats.Active++
		}

		items = append(items, cronJobRow{
			Namespace:    item.Namespace,
			Name:         item.Name,
			Schedule:     item.Spec.Schedule,
			Status:       status,
			Suspend:      yesNo(item.Spec.Suspend != nil && *item.Spec.Suspend),
			Active:       len(item.Status.Active),
			LastSchedule: formatOptionalAge(item.Status.LastScheduleTime),
			Age:          formatAge(item.CreationTimestamp),
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

func replicaSetStatus(item appsv1.ReplicaSet) string {
	desired := replicasOrZero(item.Spec.Replicas)
	switch {
	case desired == 0 && item.Status.Replicas == 0:
		return "Scaled down"
	case desired == 0:
		return "Updating"
	case item.Status.ReadyReplicas == desired && item.Status.AvailableReplicas == desired:
		return "Healthy"
	case item.Status.ReadyReplicas == 0:
		return "Degraded"
	default:
		return "Updating"
	}
}

func daemonSetStatus(item appsv1.DaemonSet) string {
	desired := item.Status.DesiredNumberScheduled
	switch {
	case desired == 0 && item.Status.CurrentNumberScheduled == 0:
		return "Scaled down"
	case desired == 0:
		return "Updating"
	case item.Status.NumberReady == desired && item.Status.UpdatedNumberScheduled == desired && item.Status.NumberMisscheduled == 0:
		return "Healthy"
	case item.Status.NumberReady == 0 || item.Status.NumberMisscheduled > 0:
		return "Degraded"
	default:
		return "Updating"
	}
}

func statefulSetStatus(item appsv1.StatefulSet) string {
	desired := desiredReplicas(item.Spec.Replicas)
	switch {
	case desired == 0 && item.Status.Replicas == 0:
		return "Scaled down"
	case item.Status.ReadyReplicas == desired && item.Status.UpdatedReplicas == desired:
		return "Healthy"
	case item.Status.ReadyReplicas == 0:
		return "Pending"
	case item.Status.UpdatedReplicas < item.Status.Replicas || item.Status.ReadyReplicas < desired:
		return "Updating"
	default:
		return "Degraded"
	}
}

func jobStatus(item batchv1.Job) string {
	switch {
	case item.Status.Failed > 0:
		return "Failed"
	case item.Status.Succeeded > 0:
		return "Succeeded"
	case item.Status.Active > 0:
		return "Running"
	default:
		return "Pending"
	}
}

func cronJobStatus(item batchv1.CronJob) string {
	switch {
	case item.Spec.Suspend != nil && *item.Spec.Suspend:
		return "Suspended"
	case len(item.Status.Active) > 0:
		return "Running"
	case item.Status.LastScheduleTime == nil:
		return "Pending"
	default:
		return "Scheduled"
	}
}

func desiredCompletions(value *int32) int32 {
	if value == nil {
		return 1
	}

	return *value
}

func jobDuration(item batchv1.Job) string {
	if item.Status.StartTime == nil {
		return "-"
	}

	if item.Status.CompletionTime != nil {
		return formatDuration(item.Status.CompletionTime.Sub(item.Status.StartTime.Time))
	}

	return formatDuration(time.Since(item.Status.StartTime.Time))
}

func formatOptionalAge(timestamp *metav1.Time) string {
	if timestamp == nil || timestamp.IsZero() {
		return "-"
	}

	return formatAge(*timestamp)
}

func yesNo(value bool) string {
	if value {
		return "Yes"
	}

	return "No"
}

func replicasOrZero(value *int32) int32 {
	if value == nil {
		return 0
	}

	return *value
}
