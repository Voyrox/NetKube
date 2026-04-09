package api

import (
	eventapi "netkube/adapters/api/cluster/v1/event"
	leaseapi "netkube/adapters/api/cluster/v1/lease"
	namespaceapi "netkube/adapters/api/cluster/v1/namespace"
	nodeapi "netkube/adapters/api/cluster/v1/node"
	configmapapi "netkube/adapters/api/configuration/v1/configmap"
	hpaapi "netkube/adapters/api/configuration/v1/hpa"
	limitrangeapi "netkube/adapters/api/configuration/v1/limitrange"
	poddisruptionbudgetapi "netkube/adapters/api/configuration/v1/poddisruptionbudget"
	resourcequotaapi "netkube/adapters/api/configuration/v1/resourcequota"
	secretapi "netkube/adapters/api/configuration/v1/secret"
	ingressapi "netkube/adapters/api/networking/v1/ingress"
	serviceapi "netkube/adapters/api/networking/v1/service"
	clusteroverview "netkube/adapters/api/overview/v1/cluster"
	workloadsoverview "netkube/adapters/api/overview/v1/workloads"
	csidriverapi "netkube/adapters/api/storage/v1/csidriver"
	csinodeapi "netkube/adapters/api/storage/v1/csinode"
	persistentvolumeapi "netkube/adapters/api/storage/v1/persistentvolume"
	persistentvolumeclaimapi "netkube/adapters/api/storage/v1/persistentvolumeclaim"
	storageclassapi "netkube/adapters/api/storage/v1/storageclass"
	volumeattachmentapi "netkube/adapters/api/storage/v1/volumeattachment"
	volumeattributeclassapi "netkube/adapters/api/storage/v1/volumeattributeclass"
	cronjobapi "netkube/adapters/api/workloads/v1/cronjob"
	daemonsetapi "netkube/adapters/api/workloads/v1/daemonset"
	deploymentapi "netkube/adapters/api/workloads/v1/deployment"
	jobapi "netkube/adapters/api/workloads/v1/job"
	podapi "netkube/adapters/api/workloads/v1/pod"
	replicasetapi "netkube/adapters/api/workloads/v1/replicaset"
	statefulsetapi "netkube/adapters/api/workloads/v1/statefulset"
)

var (
	ClusterOverviewHandler               = clusteroverview.Handler
	ClusterNodesHandler                  = nodeapi.ListHandler
	ClusterNodeDetailHandler             = nodeapi.DetailHandler
	ClusterNodeYAMLHandler               = nodeapi.YAMLHandler
	ClusterEventDetailHandler            = eventapi.DetailHandler
	ClusterNamespacesHandler             = namespaceapi.ListHandler
	ClusterLeasesHandler                 = leaseapi.ListHandler
	ClusterNamespaceYAMLHandler          = namespaceapi.YAMLHandler
	ClusterLeaseYAMLHandler              = leaseapi.YAMLHandler
	NetworkingServicesHandler            = serviceapi.ListHandler
	NetworkingIngressHandler             = ingressapi.ListHandler
	ClusterSecretsHandler                = secretapi.ListHandler
	ClusterSecretDataHandler             = secretapi.DataHandler
	ClusterConfigMapsHandler             = configmapapi.ListHandler
	ClusterHPAHandler                    = hpaapi.ListHandler
	ClusterLimitRangesHandler            = limitrangeapi.ListHandler
	ClusterResourceQuotasHandler         = resourcequotaapi.ListHandler
	ClusterPodDisruptionBudgetsHandler   = poddisruptionbudgetapi.ListHandler
	ClusterPersistentVolumesHandler      = persistentvolumeapi.ListHandler
	ClusterPersistentVolumeClaimsHandler = persistentvolumeclaimapi.ListHandler
	ClusterVolumeAttachmentsHandler      = volumeattachmentapi.ListHandler
	ClusterCSINodesHandler               = csinodeapi.ListHandler
	ClusterCSIDriversHandler             = csidriverapi.ListHandler
	ClusterStorageClassesHandler         = storageclassapi.ListHandler
	ClusterVolumeAttributeClassesHandler = volumeattributeclassapi.ListHandler
	WorkloadsOverviewHandler             = workloadsoverview.Handler
	PodsHandler                          = podapi.ListHandler
	CreatePodHandler                     = podapi.CreateHandler
	DeletePodHandler                     = podapi.DeleteHandler
	PodDetailHandler                     = podapi.DetailHandler
	PodLogsHandler                       = podapi.LogsHandler
	PodEventsHandler                     = podapi.EventsHandler
	PodYAMLHandler                       = podapi.YAMLHandler
	DeploymentsHandler                   = deploymentapi.ListHandler
	CreateDeploymentHandler              = deploymentapi.CreateHandler
	DeleteDeploymentHandler              = deploymentapi.DeleteHandler
	DeploymentDetailHandler              = deploymentapi.DetailHandler
	DeploymentEventsHandler              = deploymentapi.EventsHandler
	DeploymentYAMLHandler                = deploymentapi.YAMLHandler
	ReplicaSetsHandler                   = replicasetapi.ListHandler
	DaemonSetsHandler                    = daemonsetapi.ListHandler
	StatefulSetsHandler                  = statefulsetapi.ListHandler
	JobsHandler                          = jobapi.ListHandler
	CronJobsHandler                      = cronjobapi.ListHandler
)
