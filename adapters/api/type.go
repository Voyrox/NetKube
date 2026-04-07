package api

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

type nodeConditionRow struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type nodeListItem struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Role       string `json:"role"`
	Version    string `json:"version"`
	InternalIP string `json:"internalIP"`
	Age        string `json:"age"`
}

type nodeListResponse struct {
	Meta  pageMeta       `json:"meta"`
	Items []nodeListItem `json:"items"`
	Count int            `json:"count"`
}

type namespaceListItem struct {
	Name  string `json:"name"`
	Phase string `json:"phase"`
	Age   string `json:"age"`
}

type namespaceListResponse struct {
	Meta  pageMeta            `json:"meta"`
	Items []namespaceListItem `json:"items"`
	Count int                 `json:"count"`
}

type leaseListItem struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Holder    string `json:"holder"`
	LastRenew string `json:"lastRenew"`
	Age       string `json:"age"`
}

type leaseListResponse struct {
	Meta  pageMeta        `json:"meta"`
	Items []leaseListItem `json:"items"`
	Count int             `json:"count"`
}

type serviceListItem struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	ExternalIP string `json:"externalIP"`
	Ports      string `json:"ports"`
	Selector   string `json:"selector"`
	Age        string `json:"age"`
}

type serviceListResponse struct {
	Meta  pageMeta          `json:"meta"`
	Items []serviceListItem `json:"items"`
	Count int               `json:"count"`
}

type nodeEventRow struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Age     string `json:"age"`
	Type    string `json:"type"`
}

type nodeDetail struct {
	Name              string             `json:"name"`
	Status            string             `json:"status"`
	Role              string             `json:"role"`
	KubeletVersion    string             `json:"kubeletVersion"`
	ContainerRuntime  string             `json:"containerRuntime"`
	OSKernel          string             `json:"osKernel"`
	Architecture      string             `json:"architecture"`
	InternalIP        string             `json:"internalIP"`
	PodCIDR           string             `json:"podCIDR"`
	AllocatableCPU    string             `json:"allocatableCPU"`
	AllocatableMemory string             `json:"allocatableMemory"`
	AllocatablePods   string             `json:"allocatablePods"`
	AllocatableStore  string             `json:"allocatableStore"`
	CPUPressure       string             `json:"cpuPressure"`
	MemoryPressure    string             `json:"memoryPressure"`
	DiskPressure      string             `json:"diskPressure"`
	Labels            map[string]string  `json:"labels"`
	Conditions        []nodeConditionRow `json:"conditions"`
	Timeline          []nodeEventRow     `json:"timeline"`
}

type nodeDetailResponse struct {
	Meta pageMeta   `json:"meta"`
	Item nodeDetail `json:"item"`
}

type clusterEventDetail struct {
	Title          string            `json:"title"`
	Type           string            `json:"type"`
	Namespace      string            `json:"namespace"`
	Reason         string            `json:"reason"`
	InvolvedObject string            `json:"involvedObject"`
	Kind           string            `json:"kind"`
	Name           string            `json:"name"`
	Source         string            `json:"source"`
	FirstSeen      string            `json:"firstSeen"`
	LastSeen       string            `json:"lastSeen"`
	Count          int32             `json:"count"`
	Node           string            `json:"node"`
	Message        string            `json:"message"`
	Timeline       []nodeEventRow    `json:"timeline"`
	Annotations    map[string]string `json:"annotations"`
}

type clusterEventDetailResponse struct {
	Meta pageMeta           `json:"meta"`
	Item clusterEventDetail `json:"item"`
}
