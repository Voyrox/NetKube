package api

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

type deploymentRow struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Ready     string `json:"ready"`
	Status    string `json:"status"`
	Desired   int32  `json:"desired"`
	Updated   int32  `json:"updated"`
	Available int32  `json:"available"`
	Age       string `json:"age"`
}

type deploymentsResponse struct {
	Meta  pageMeta         `json:"meta"`
	Items []deploymentRow  `json:"items"`
	Count int              `json:"count"`
	Error string           `json:"error,omitempty"`
	Stats deploymentsStats `json:"stats"`
}

type deploymentsStats struct {
	Healthy int `json:"healthy"`
	Warning int `json:"warning"`
	Pending int `json:"pending"`
	Total   int `json:"total"`
}

type deploymentConditionRow struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type deploymentReplicaSetRow struct {
	Name    string `json:"name"`
	Ready   string `json:"ready"`
	Desired int32  `json:"desired"`
	Age     string `json:"age"`
}

type deploymentPodRow struct {
	Name   string `json:"name"`
	Ready  string `json:"ready"`
	Status string `json:"status"`
	Node   string `json:"node"`
	Age    string `json:"age"`
}

type deploymentDetail struct {
	Namespace   string                    `json:"namespace"`
	Name        string                    `json:"name"`
	Status      string                    `json:"status"`
	Ready       string                    `json:"ready"`
	Desired     int32                     `json:"desired"`
	Updated     int32                     `json:"updated"`
	Available   int32                     `json:"available"`
	Unavailable int32                     `json:"unavailable"`
	Age         string                    `json:"age"`
	Strategy    string                    `json:"strategy"`
	Selector    string                    `json:"selector"`
	Conditions  []deploymentConditionRow  `json:"conditions"`
	ReplicaSets []deploymentReplicaSetRow `json:"replicaSets"`
	Pods        []deploymentPodRow        `json:"pods"`
	Labels      map[string]string         `json:"labels"`
	Annotations map[string]string         `json:"annotations"`
}

type deploymentDetailResponse struct {
	Meta pageMeta         `json:"meta"`
	Item deploymentDetail `json:"item"`
}

type deploymentEventRow struct {
	Type    string `json:"type"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Age     string `json:"age"`
}

type deploymentEventsResponse struct {
	Meta  pageMeta             `json:"meta"`
	Items []deploymentEventRow `json:"items"`
}

type deploymentYAMLResponse struct {
	Meta    pageMeta `json:"meta"`
	Content string   `json:"content"`
}

func DeploymentsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := selectedNamespace(c)
	items, stats, err := GetDeployments(cluster.Clientset, namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deploymentsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
		Count: len(items),
		Stats: stats,
	})
}

func DeploymentDetailHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "deployment name and namespace are required"})
		return
	}

	item, err := GetDeploymentDetail(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deploymentDetailResponse{
		Meta: pageMetaFromCluster(cluster, namespace),
		Item: item,
	})
}

func DeploymentEventsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "deployment name and namespace are required"})
		return
	}

	items, err := GetDeploymentEvents(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deploymentEventsResponse{
		Meta:  pageMetaFromCluster(cluster, namespace),
		Items: items,
	})
}

func DeploymentYAMLHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "deployment name and namespace are required"})
		return
	}

	content, err := GetDeploymentYAML(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deploymentYAMLResponse{
		Meta:    pageMetaFromCluster(cluster, namespace),
		Content: content,
	})
}

func GetDeployments(clientset *kubernetes.Clientset, namespace string) ([]deploymentRow, deploymentsStats, error) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, deploymentsStats{}, err
	}

	items := make([]deploymentRow, 0, len(deployments.Items))
	stats := deploymentsStats{Total: len(deployments.Items)}

	for _, item := range deployments.Items {
		status := deploymentStatus(item)
		switch status {
		case "Healthy":
			stats.Healthy++
		case "Pending":
			stats.Pending++
		default:
			stats.Warning++
		}

		items = append(items, deploymentRow{
			Namespace: item.Namespace,
			Name:      item.Name,
			Ready:     readyRatio(item.Status.ReadyReplicas, item.Status.Replicas),
			Status:    status,
			Desired:   item.Status.Replicas,
			Updated:   item.Status.UpdatedReplicas,
			Available: item.Status.AvailableReplicas,
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

func GetDeploymentDetail(clientset *kubernetes.Clientset, namespace, name string) (deploymentDetail, error) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return deploymentDetail{}, err
	}

	selector := labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels)
	selectorString := selector.String()

	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: selectorString})
	if err != nil {
		return deploymentDetail{}, err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: selectorString})
	if err != nil {
		return deploymentDetail{}, err
	}

	conditions := make([]deploymentConditionRow, 0, len(deployment.Status.Conditions))
	for _, condition := range deployment.Status.Conditions {
		conditions = append(conditions, deploymentConditionRow{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  fallback(condition.Reason),
			Message: fallback(condition.Message),
		})
	}

	replicaSetItems := make([]deploymentReplicaSetRow, 0, len(replicaSets.Items))
	for _, item := range replicaSets.Items {
		replicaSetItems = append(replicaSetItems, deploymentReplicaSetRow{
			Name:    item.Name,
			Ready:   readyRatio(item.Status.ReadyReplicas, item.Status.Replicas),
			Desired: item.Status.Replicas,
			Age:     formatAge(item.CreationTimestamp),
		})
	}
	sort.Slice(replicaSetItems, func(i, j int) bool { return replicaSetItems[i].Name < replicaSetItems[j].Name })

	podItems := make([]deploymentPodRow, 0, len(pods.Items))
	for _, item := range pods.Items {
		podItems = append(podItems, deploymentPodRow{
			Name:   item.Name,
			Ready:  readyRatio(readyContainers(item), int32(len(item.Status.ContainerStatuses))),
			Status: podStatus(item),
			Node:   fallback(item.Spec.NodeName),
			Age:    formatAge(item.CreationTimestamp),
		})
	}
	sort.Slice(podItems, func(i, j int) bool { return podItems[i].Name < podItems[j].Name })

	return deploymentDetail{
		Namespace:   deployment.Namespace,
		Name:        deployment.Name,
		Status:      deploymentStatus(*deployment),
		Ready:       readyRatio(deployment.Status.ReadyReplicas, deployment.Status.Replicas),
		Desired:     deployment.Status.Replicas,
		Updated:     deployment.Status.UpdatedReplicas,
		Available:   deployment.Status.AvailableReplicas,
		Unavailable: deployment.Status.UnavailableReplicas,
		Age:         formatAge(deployment.CreationTimestamp),
		Strategy:    string(deployment.Spec.Strategy.Type),
		Selector:    fallback(selectorString),
		Conditions:  conditions,
		ReplicaSets: replicaSetItems,
		Pods:        podItems,
		Labels:      cloneStringMap(deployment.Labels),
		Annotations: cloneStringMap(deployment.Annotations),
	}, nil
}

func GetDeploymentEvents(clientset *kubernetes.Clientset, namespace, name string) ([]deploymentEventRow, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{
		FieldSelector: "involvedObject.kind=Deployment,involvedObject.name=" + name,
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].CreationTimestamp.Time.After(events.Items[j].CreationTimestamp.Time)
	})

	items := make([]deploymentEventRow, 0, len(events.Items))
	for _, item := range events.Items {
		items = append(items, deploymentEventRow{
			Type:    fallback(item.Type),
			Reason:  fallback(item.Reason),
			Message: fallback(item.Message),
			Age:     formatAge(item.CreationTimestamp),
		})
	}

	return items, nil
}

func GetDeploymentYAML(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
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

func deploymentStatus(item appsv1.Deployment) string {
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
