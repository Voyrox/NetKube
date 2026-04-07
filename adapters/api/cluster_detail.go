package api

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func ClusterNodeDetailHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	item, err := GetNodeDetail(cluster.Clientset, strings.TrimSpace(c.Query("name")))
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodeDetailResponse{
		Meta: pageMetaFromCluster(cluster, ""),
		Item: item,
	})
}

func ClusterNodesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetNodeList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodeListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterNamespacesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetNamespaceList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespaceListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterLeasesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetLeaseList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaseListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func NetworkingServicesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetServiceList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, serviceListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterNamespaceYAMLHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "namespace name is required"})
		return
	}

	namespace, err := cluster.Clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	content, err := yaml.Marshal(namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": pageMetaFromCluster(cluster, ""), "content": string(content)})
}

func ClusterLeaseYAMLHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	namespace := strings.TrimSpace(c.Query("namespace"))
	name := strings.TrimSpace(c.Query("name"))
	if namespace == "" || name == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "lease namespace and name are required"})
		return
	}

	lease, err := cluster.Clientset.CoordinationV1().Leases(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	content, err := yaml.Marshal(lease)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meta": pageMetaFromCluster(cluster, namespace), "content": string(content)})
}

func ClusterEventDetailHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	item, err := GetClusterEventDetail(cluster.Clientset, strings.TrimSpace(c.Query("namespace")), strings.TrimSpace(c.Query("name")), strings.TrimSpace(c.Query("reason")), strings.TrimSpace(c.Query("kind")))
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, clusterEventDetailResponse{
		Meta: pageMetaFromCluster(cluster, item.Namespace),
		Item: item,
	})
}

func GetNodeDetail(clientset *kubernetes.Clientset, name string) (nodeDetail, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nodeDetail{}, err
	}
	if len(nodes.Items) == 0 {
		return nodeDetail{}, nil
	}

	sort.Slice(nodes.Items, func(i, j int) bool { return nodes.Items[i].Name < nodes.Items[j].Name })
	node := nodes.Items[0]
	if name != "" {
		for _, item := range nodes.Items {
			if item.Name == name {
				node = item
				break
			}
		}
	}

	conditions := make([]nodeConditionRow, 0, len(node.Status.Conditions))
	var cpuPressure, memoryPressure, diskPressure string
	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, nodeConditionRow{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  fallback(condition.Reason),
			Message: fallback(condition.Message),
		})
		switch condition.Type {
		case corev1.NodePIDPressure:
			if cpuPressure == "" {
				cpuPressure = boolStateLabel(condition.Status == corev1.ConditionFalse)
			}
		case corev1.NodeMemoryPressure:
			memoryPressure = boolStateLabel(condition.Status == corev1.ConditionFalse)
		case corev1.NodeDiskPressure:
			diskPressure = boolStateLabel(condition.Status == corev1.ConditionFalse)
		}
	}

	events, _ := clientset.CoreV1().Events("").List(context.Background(), metav1.ListOptions{
		FieldSelector: "involvedObject.kind=Node,involvedObject.name=" + node.Name,
	})
	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].CreationTimestamp.Time.After(events.Items[j].CreationTimestamp.Time)
	})
	timeline := make([]nodeEventRow, 0, minInt(len(events.Items), 5))
	for _, item := range events.Items {
		timeline = append(timeline, nodeEventRow{Title: fallback(item.Reason), Message: fallback(item.Message), Age: formatAge(item.CreationTimestamp), Type: fallback(item.Type)})
		if len(timeline) == 5 {
			break
		}
	}

	allocatableStore := int64(0)
	if quantity, ok := node.Status.Allocatable[corev1.ResourceEphemeralStorage]; ok {
		allocatableStore = quantity.Value()
	}

	return nodeDetail{
		Name:              node.Name,
		Status:            nodeStatus(node),
		Role:              nodeRole(node),
		KubeletVersion:    fallback(node.Status.NodeInfo.KubeletVersion),
		ContainerRuntime:  fallback(node.Status.NodeInfo.ContainerRuntimeVersion),
		OSKernel:          strings.TrimSpace(strings.Join([]string{fallback(node.Status.NodeInfo.OSImage), fallback(node.Status.NodeInfo.KernelVersion)}, " / ")),
		Architecture:      fallback(node.Status.NodeInfo.Architecture),
		InternalIP:        nodeAddress(node, corev1.NodeInternalIP),
		PodCIDR:           fallback(node.Spec.PodCIDR),
		AllocatableCPU:    formatCPU(node.Status.Allocatable.Cpu().MilliValue()),
		AllocatableMemory: formatBytes(node.Status.Allocatable.Memory().Value()),
		AllocatablePods:   fallback(node.Status.Allocatable.Pods().String()),
		AllocatableStore:  formatBytes(allocatableStore),
		CPUPressure:       fallback(cpuPressure),
		MemoryPressure:    fallback(memoryPressure),
		DiskPressure:      fallback(diskPressure),
		Labels:            cloneStringMap(node.Labels),
		Conditions:        conditions,
		Timeline:          timeline,
	}, nil
}

func GetNodeList(clientset *kubernetes.Clientset) ([]nodeListItem, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]nodeListItem, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		items = append(items, nodeListItem{
			Name:       node.Name,
			Status:     nodeStatus(node),
			Role:       nodeRole(node),
			Version:    fallback(node.Status.NodeInfo.KubeletVersion),
			InternalIP: nodeAddress(node, corev1.NodeInternalIP),
			Age:        formatAge(node.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func GetNamespaceList(clientset *kubernetes.Clientset) ([]namespaceListItem, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]namespaceListItem, 0, len(namespaces.Items))
	for _, item := range namespaces.Items {
		items = append(items, namespaceListItem{
			Name:  item.Name,
			Phase: fallback(string(item.Status.Phase)),
			Age:   formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func GetLeaseList(clientset *kubernetes.Clientset) ([]leaseListItem, error) {
	leases, err := clientset.CoordinationV1().Leases("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]leaseListItem, 0, len(leases.Items))
	for _, item := range leases.Items {
		items = append(items, leaseListItem{
			Namespace: item.Namespace,
			Name:      item.Name,
			Holder:    leaseHolder(item),
			LastRenew: leaseRenew(item),
			Age:       formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})
	return items, nil
}

func GetServiceList(clientset *kubernetes.Clientset) ([]serviceListItem, error) {
	services, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]serviceListItem, 0, len(services.Items))
	for _, item := range services.Items {
		items = append(items, serviceListItem{
			Namespace:  item.Namespace,
			Name:       item.Name,
			Type:       fallback(string(item.Spec.Type)),
			ExternalIP: serviceExternalIP(item),
			Ports:      servicePorts(item),
			Selector:   mapToSelector(item.Spec.Selector),
			Age:        formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Namespace == items[j].Namespace {
			return items[i].Name < items[j].Name
		}
		return items[i].Namespace < items[j].Namespace
	})

	return items, nil
}

func GetClusterEventDetail(clientset *kubernetes.Clientset, namespace, name, reason, kind string) (clusterEventDetail, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return clusterEventDetail{}, err
	}
	sort.Slice(events.Items, func(i, j int) bool {
		return events.Items[i].CreationTimestamp.Time.After(events.Items[j].CreationTimestamp.Time)
	})
	if len(events.Items) == 0 {
		return clusterEventDetail{}, nil
	}
	selected := events.Items[0]
	for _, item := range events.Items {
		if (namespace == "" || item.Namespace == namespace) && (name == "" || item.InvolvedObject.Name == name) && (reason == "" || item.Reason == reason) && (kind == "" || item.InvolvedObject.Kind == kind) {
			selected = item
			break
		}
	}

	related := make([]nodeEventRow, 0, 5)
	for _, item := range events.Items {
		if item.InvolvedObject.Kind == selected.InvolvedObject.Kind && item.InvolvedObject.Name == selected.InvolvedObject.Name {
			related = append(related, nodeEventRow{Title: fallback(item.Reason), Message: fallback(item.Message), Age: formatAge(item.CreationTimestamp), Type: fallback(item.Type)})
			if len(related) == 5 {
				break
			}
		}
	}

	return clusterEventDetail{
		Title:          fallback(selected.Message),
		Type:           fallback(selected.Type),
		Namespace:      fallback(selected.Namespace),
		Reason:         fallback(selected.Reason),
		InvolvedObject: strings.TrimSpace(strings.Join([]string{fallback(selected.InvolvedObject.Kind), fallback(selected.InvolvedObject.Name)}, " / ")),
		Kind:           fallback(selected.InvolvedObject.Kind),
		Name:           fallback(selected.InvolvedObject.Name),
		Source:         strings.TrimSpace(strings.Join([]string{fallback(selected.Source.Component), fallback(selected.Source.Host)}, " / ")),
		FirstSeen:      formatAge(selected.CreationTimestamp),
		LastSeen:       formatAge(selected.CreationTimestamp),
		Count:          selected.Count,
		Node:           fallback(selected.Source.Host),
		Message:        fallback(selected.Message),
		Timeline:       related,
		Annotations:    cloneStringMap(selected.Annotations),
	}, nil
}

func nodeRole(node corev1.Node) string {
	for key := range node.Labels {
		if strings.HasPrefix(key, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(key, "node-role.kubernetes.io/")
			if role == "" {
				return "Control plane"
			}
			return strings.ToUpper(role[:1]) + role[1:]
		}
	}
	return "Worker"
}

func nodeStatus(node corev1.Node) string {
	if nodeReady(node) {
		return "Ready"
	}
	if node.Spec.Unschedulable {
		return "Cordoned"
	}
	return "Not ready"
}

func nodeAddress(node corev1.Node, addressType corev1.NodeAddressType) string {
	for _, address := range node.Status.Addresses {
		if address.Type == addressType {
			return fallback(address.Address)
		}
	}
	return "-"
}

func boolStateLabel(ok bool) string {
	if ok {
		return "Clear"
	}
	return "Detected"
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func leaseHolder(lease coordinationv1.Lease) string {
	if lease.Spec.HolderIdentity == nil || strings.TrimSpace(*lease.Spec.HolderIdentity) == "" {
		return "-"
	}
	return *lease.Spec.HolderIdentity
}

func leaseRenew(lease coordinationv1.Lease) string {
	if lease.Spec.RenewTime == nil || lease.Spec.RenewTime.IsZero() {
		return "-"
	}
	return formatAge(metav1.Time{Time: lease.Spec.RenewTime.Time})
}

func serviceExternalIP(service corev1.Service) string {
	if len(service.Spec.ExternalIPs) > 0 {
		return strings.Join(service.Spec.ExternalIPs, ", ")
	}
	if len(service.Status.LoadBalancer.Ingress) > 0 {
		values := make([]string, 0, len(service.Status.LoadBalancer.Ingress))
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			if ingress.IP != "" {
				values = append(values, ingress.IP)
				continue
			}
			if ingress.Hostname != "" {
				values = append(values, ingress.Hostname)
			}
		}
		if len(values) > 0 {
			return strings.Join(values, ", ")
		}
	}
	if service.Spec.ClusterIP != "" && service.Spec.ClusterIP != "None" {
		return service.Spec.ClusterIP
	}
	return "-"
}

func servicePorts(service corev1.Service) string {
	if len(service.Spec.Ports) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(service.Spec.Ports))
	for _, port := range service.Spec.Ports {
		part := string(port.Protocol) + "/" + int32String(port.Port)
		if port.TargetPort.String() != "" {
			part += " -> " + port.TargetPort.String()
		}
		if strings.TrimSpace(port.Name) != "" {
			part += " (" + port.Name + ")"
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, ", ")
}

func mapToSelector(values map[string]string) string {
	if len(values) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(values))
	for key, value := range values {
		parts = append(parts, key+"="+value)
	}
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}

func ClusterNodeYAMLHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "node name is required"})
		return
	}
	node, err := cluster.Clientset.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}
	content, err := yaml.Marshal(node)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"meta": pageMetaFromCluster(cluster, ""), "content": string(content)})
}
