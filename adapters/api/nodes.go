package api

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type nodesSummary struct {
	Total    int    `json:"total"`
	Ready    int    `json:"ready"`
	Pending  int    `json:"pending"`
	NotReady int    `json:"notReady"`
	Status   string `json:"status"`
}

func GetNodes(clientset *kubernetes.Clientset) (nodesSummary, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nodesSummary{}, err
	}

	summary := nodesSummary{Total: len(nodes.Items)}
	for _, node := range nodes.Items {
		if nodeReady(node) {
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

func nodeReady(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}

	return false
}
