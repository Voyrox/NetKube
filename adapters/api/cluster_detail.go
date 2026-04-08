package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

func NetworkingIngressHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetIngressList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ingressListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterSecretsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetSecretList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, secretListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterConfigMapsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetConfigMapList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, configMapListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterHPAHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetHPAList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, hpaListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterLimitRangesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetLimitRangeList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, limitRangeListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterResourceQuotasHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetResourceQuotaList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resourceQuotaListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterPodDisruptionBudgetsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetPodDisruptionBudgetList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, podDisruptionBudgetListResponse{
		Meta:  pageMetaFromCluster(cluster, ""),
		Items: items,
		Count: len(items),
	})
}

func ClusterPersistentVolumesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetPersistentVolumeList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, persistentVolumeListResponse{Meta: pageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func ClusterPersistentVolumeClaimsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetPersistentVolumeClaimList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, persistentVolumeClaimListResponse{Meta: pageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func ClusterVolumeAttachmentsHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetVolumeAttachmentList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, volumeAttachmentListResponse{Meta: pageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func ClusterCSINodesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetCSINodeList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, csiNodeListResponse{Meta: pageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func ClusterCSIDriversHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetCSIDriverList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, csiDriverListResponse{Meta: pageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func ClusterStorageClassesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetStorageClassList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, storageClassListResponse{Meta: pageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func ClusterVolumeAttributeClassesHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	items, err := GetVolumeAttributeClassList(cluster.Clientset)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, volumeAttributeClassListResponse{Meta: pageMetaFromCluster(cluster, ""), Items: items, Count: len(items)})
}

func ClusterSecretDataHandler(c *gin.Context) {
	cluster, ok := resolveClusterRequest(c)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Query("name"))
	namespace := strings.TrimSpace(c.Query("namespace"))
	if name == "" || namespace == "" {
		c.JSON(http.StatusBadRequest, apiError{Error: "secret name and namespace are required"})
		return
	}

	content, err := GetSecretData(cluster.Clientset, namespace, name)
	if err != nil {
		c.JSON(http.StatusBadGateway, apiError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, secretDataResponse{
		Meta:      pageMetaFromCluster(cluster, namespace),
		Namespace: namespace,
		Name:      name,
		Content:   content,
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

func GetIngressList(clientset *kubernetes.Clientset) ([]ingressListItem, error) {
	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]ingressListItem, 0, len(ingresses.Items))
	for _, item := range ingresses.Items {
		items = append(items, ingressListItem{
			Namespace: item.Namespace,
			Name:      item.Name,
			Class:     ingressClass(item),
			Hosts:     ingressHosts(item),
			Address:   ingressAddress(item),
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

func GetSecretList(clientset *kubernetes.Clientset) ([]secretListItem, error) {
	secrets, err := clientset.CoreV1().Secrets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]secretListItem, 0, len(secrets.Items))
	for _, item := range secrets.Items {
		items = append(items, secretListItem{
			Namespace: item.Namespace,
			Name:      item.Name,
			Type:      fallback(string(item.Type)),
			Data:      len(item.Data),
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

func GetConfigMapList(clientset *kubernetes.Clientset) ([]configMapListItem, error) {
	configMaps, err := clientset.CoreV1().ConfigMaps("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]configMapListItem, 0, len(configMaps.Items))
	for _, item := range configMaps.Items {
		items = append(items, configMapListItem{
			Namespace: item.Namespace,
			Name:      item.Name,
			Data:      len(item.Data) + len(item.BinaryData),
			Immutable: yesNo(item.Immutable != nil && *item.Immutable),
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

func GetSecretData(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	content := make(map[string]string, len(secret.Data))
	for key, value := range secret.Data {
		content[key] = secretDisplayValue(value)
	}

	rendered, err := yaml.Marshal(content)
	if err != nil {
		return "", err
	}

	return string(rendered), nil
}

func GetHPAList(clientset *kubernetes.Clientset) ([]hpaListItem, error) {
	hpas, err := clientset.AutoscalingV2().HorizontalPodAutoscalers("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]hpaListItem, 0, len(hpas.Items))
	for _, item := range hpas.Items {
		items = append(items, hpaListItem{
			Namespace: item.Namespace,
			Name:      item.Name,
			Target:    hpaTarget(item),
			Min:       int32PointerString(item.Spec.MinReplicas),
			Max:       item.Spec.MaxReplicas,
			Current:   int32String(item.Status.CurrentReplicas),
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

func GetLimitRangeList(clientset *kubernetes.Clientset) ([]limitRangeListItem, error) {
	limitRanges, err := clientset.CoreV1().LimitRanges("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]limitRangeListItem, 0, len(limitRanges.Items))
	for _, item := range limitRanges.Items {
		items = append(items, limitRangeListItem{
			Namespace: item.Namespace,
			Name:      item.Name,
			Limits:    len(item.Spec.Limits),
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

func GetResourceQuotaList(clientset *kubernetes.Clientset) ([]resourceQuotaListItem, error) {
	resourceQuotas, err := clientset.CoreV1().ResourceQuotas("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]resourceQuotaListItem, 0, len(resourceQuotas.Items))
	for _, item := range resourceQuotas.Items {
		items = append(items, resourceQuotaListItem{
			Namespace: item.Namespace,
			Name:      item.Name,
			Scopes:    len(item.Spec.Scopes),
			Hard:      len(item.Status.Hard),
			Used:      len(item.Status.Used),
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

func GetPodDisruptionBudgetList(clientset *kubernetes.Clientset) ([]podDisruptionBudgetListItem, error) {
	pdbs, err := clientset.PolicyV1().PodDisruptionBudgets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]podDisruptionBudgetListItem, 0, len(pdbs.Items))
	for _, item := range pdbs.Items {
		items = append(items, podDisruptionBudgetListItem{
			Namespace:      item.Namespace,
			Name:           item.Name,
			MinAvailable:   intOrStringValue(item.Spec.MinAvailable),
			MaxUnavailable: intOrStringValue(item.Spec.MaxUnavailable),
			Allowed:        item.Status.DisruptionsAllowed,
			Age:            formatAge(item.CreationTimestamp),
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

func GetPersistentVolumeList(clientset *kubernetes.Clientset) ([]persistentVolumeListItem, error) {
	volumes, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]persistentVolumeListItem, 0, len(volumes.Items))
	for _, item := range volumes.Items {
		items = append(items, persistentVolumeListItem{
			Name:          item.Name,
			Status:        fallback(string(item.Status.Phase)),
			Capacity:      resourceListValue(item.Spec.Capacity, corev1.ResourceStorage),
			AccessModes:   pvAccessModes(item.Spec.AccessModes),
			ReclaimPolicy: fallback(string(item.Spec.PersistentVolumeReclaimPolicy)),
			Claim:         pvClaim(item.Spec.ClaimRef),
			StorageClass:  fallback(item.Spec.StorageClassName),
			Age:           formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func GetPersistentVolumeClaimList(clientset *kubernetes.Clientset) ([]persistentVolumeClaimListItem, error) {
	claims, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]persistentVolumeClaimListItem, 0, len(claims.Items))
	for _, item := range claims.Items {
		items = append(items, persistentVolumeClaimListItem{
			Namespace:    item.Namespace,
			Name:         item.Name,
			Status:       fallback(string(item.Status.Phase)),
			Volume:       fallback(item.Spec.VolumeName),
			Capacity:     resourceListValue(item.Status.Capacity, corev1.ResourceStorage),
			AccessModes:  pvcAccessModes(item.Status.AccessModes, item.Spec.AccessModes),
			StorageClass: pvcStorageClass(item.Spec.StorageClassName),
			Age:          formatAge(item.CreationTimestamp),
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

func GetVolumeAttachmentList(clientset *kubernetes.Clientset) ([]volumeAttachmentListItem, error) {
	attachments, err := clientset.StorageV1().VolumeAttachments().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]volumeAttachmentListItem, 0, len(attachments.Items))
	for _, item := range attachments.Items {
		items = append(items, volumeAttachmentListItem{
			Name:             item.Name,
			Attacher:         fallback(item.Spec.Attacher),
			Node:             fallback(item.Spec.NodeName),
			PersistentVolume: fallback(stringPointer(item.Spec.Source.PersistentVolumeName)),
			Attached:         yesNo(item.Status.Attached),
			Age:              formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func GetCSINodeList(clientset *kubernetes.Clientset) ([]csiNodeListItem, error) {
	nodes, err := clientset.StorageV1().CSINodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]csiNodeListItem, 0, len(nodes.Items))
	for _, item := range nodes.Items {
		items = append(items, csiNodeListItem{
			Name:    item.Name,
			Drivers: len(item.Spec.Drivers),
			Age:     formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func GetCSIDriverList(clientset *kubernetes.Clientset) ([]csiDriverListItem, error) {
	drivers, err := clientset.StorageV1().CSIDrivers().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]csiDriverListItem, 0, len(drivers.Items))
	for _, item := range drivers.Items {
		items = append(items, csiDriverListItem{
			Name:            item.Name,
			AttachRequired:  yesNo(boolPointer(item.Spec.AttachRequired)),
			PodInfoOnMount:  yesNo(boolPointer(item.Spec.PodInfoOnMount)),
			StorageCapacity: yesNo(boolPointer(item.Spec.StorageCapacity)),
			Modes:           volumeLifecycleModes(item.Spec.VolumeLifecycleModes),
			Age:             formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func GetStorageClassList(clientset *kubernetes.Clientset) ([]storageClassListItem, error) {
	classes, err := clientset.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]storageClassListItem, 0, len(classes.Items))
	for _, item := range classes.Items {
		items = append(items, storageClassListItem{
			Name:          item.Name,
			Provisioner:   fallback(item.Provisioner),
			ReclaimPolicy: reclaimPolicy(item.ReclaimPolicy),
			BindingMode:   bindingMode(item.VolumeBindingMode),
			Default:       yesNo(isDefaultStorageClass(item)),
			Age:           formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func GetVolumeAttributeClassList(clientset *kubernetes.Clientset) ([]volumeAttributeClassListItem, error) {
	classes, err := clientset.StorageV1().VolumeAttributesClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]volumeAttributeClassListItem, 0, len(classes.Items))
	for _, item := range classes.Items {
		items = append(items, volumeAttributeClassListItem{
			Name:       item.Name,
			DriverName: fallback(item.DriverName),
			Parameters: len(item.Parameters),
			Age:        formatAge(item.CreationTimestamp),
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
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

func ingressClass(ingress networkingv1.Ingress) string {
	if ingress.Spec.IngressClassName != nil && strings.TrimSpace(*ingress.Spec.IngressClassName) != "" {
		return *ingress.Spec.IngressClassName
	}
	return "-"
}

func ingressHosts(ingress networkingv1.Ingress) string {
	if len(ingress.Spec.Rules) == 0 {
		return "-"
	}

	hosts := make([]string, 0, len(ingress.Spec.Rules))
	for _, rule := range ingress.Spec.Rules {
		if strings.TrimSpace(rule.Host) != "" {
			hosts = append(hosts, rule.Host)
		}
	}
	if len(hosts) == 0 {
		return "-"
	}
	return strings.Join(hosts, ", ")
}

func ingressAddress(ingress networkingv1.Ingress) string {
	if len(ingress.Status.LoadBalancer.Ingress) == 0 {
		return "-"
	}

	addresses := make([]string, 0, len(ingress.Status.LoadBalancer.Ingress))
	for _, entry := range ingress.Status.LoadBalancer.Ingress {
		if strings.TrimSpace(entry.IP) != "" {
			addresses = append(addresses, entry.IP)
			continue
		}
		if strings.TrimSpace(entry.Hostname) != "" {
			addresses = append(addresses, entry.Hostname)
		}
	}
	if len(addresses) == 0 {
		return "-"
	}
	return strings.Join(addresses, ", ")
}

func secretDisplayValue(value []byte) string {
	if utf8.Valid(value) {
		return string(value)
	}

	return "base64:" + base64.StdEncoding.EncodeToString(value)
}

func hpaTarget(item autoscalingv2.HorizontalPodAutoscaler) string {
	kind := strings.TrimSpace(item.Spec.ScaleTargetRef.Kind)
	name := strings.TrimSpace(item.Spec.ScaleTargetRef.Name)
	if kind == "" && name == "" {
		return "-"
	}
	return strings.TrimSpace(strings.Join([]string{kind, name}, "/"))
}

func int32PointerString(value *int32) string {
	if value == nil {
		return "-"
	}
	return int32String(*value)
}

func intOrStringValue(value *intstr.IntOrString) string {
	if value == nil {
		return "-"
	}
	return value.String()
}

func pdbPolicyVersion(_ policyv1.PodDisruptionBudget) string {
	return "policy/v1"
}

func resourceListValue(values corev1.ResourceList, key corev1.ResourceName) string {
	quantity, ok := values[key]
	if !ok {
		return "-"
	}
	return quantity.String()
}

func pvAccessModes(modes []corev1.PersistentVolumeAccessMode) string {
	if len(modes) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(modes))
	for _, mode := range modes {
		parts = append(parts, string(mode))
	}
	return strings.Join(parts, ", ")
}

func pvcAccessModes(primary, fallbackModes []corev1.PersistentVolumeAccessMode) string {
	if len(primary) > 0 {
		return pvAccessModes(primary)
	}
	return pvAccessModes(fallbackModes)
}

func pvClaim(claimRef *corev1.ObjectReference) string {
	if claimRef == nil {
		return "-"
	}
	if strings.TrimSpace(claimRef.Namespace) == "" {
		return fallback(claimRef.Name)
	}
	return strings.TrimSpace(claimRef.Namespace + "/" + claimRef.Name)
}

func pvcStorageClass(name *string) string {
	if name == nil || strings.TrimSpace(*name) == "" {
		return "-"
	}
	return *name
}

func stringPointer(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func boolPointer(value *bool) bool {
	return value != nil && *value
}

func volumeLifecycleModes(modes []storagev1.VolumeLifecycleMode) string {
	if len(modes) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(modes))
	for _, mode := range modes {
		parts = append(parts, string(mode))
	}
	return strings.Join(parts, ", ")
}

func reclaimPolicy(value *corev1.PersistentVolumeReclaimPolicy) string {
	if value == nil {
		return "-"
	}
	return string(*value)
}

func bindingMode(value *storagev1.VolumeBindingMode) string {
	if value == nil {
		return "-"
	}
	return string(*value)
}

func isDefaultStorageClass(item storagev1.StorageClass) bool {
	return item.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" || item.Annotations["storageclass.beta.kubernetes.io/is-default-class"] == "true"
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
