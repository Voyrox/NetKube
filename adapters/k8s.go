package adapters

import (
	"fmt"
	"time"

	"netkube/storage"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type ClusterClient struct {
	Context             storage.ContextItem
	Clientset           *kubernetes.Clientset
	APIExtensionsClient *apiextensionsclient.Clientset
}

func ResolveCluster(contextID string) (*ClusterClient, error) {
	contexts, err := storage.LoadContexts()
	if err != nil {
		return nil, fmt.Errorf("load contexts: %w", err)
	}

	var selected storage.ContextItem
	found := false
	for _, item := range contexts {
		if item.ID == contextID {
			selected = item
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("context %q not found", contextID)
	}

	sources, err := storage.LoadSources()
	if err != nil {
		return nil, fmt.Errorf("load sources: %w", err)
	}

	var source storage.SourceItem
	found = false
	for _, item := range sources {
		if item.ID == selected.SourceID {
			source = item
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("source %q not found for context %q", selected.SourceID, selected.ContextName)
	}

	storedPath, err := storage.StoredSourcePath(source)
	if err != nil {
		return nil, fmt.Errorf("load kubeconfig source %q: %w", source.Name, err)
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: storedPath}
	overrides := &clientcmd.ConfigOverrides{CurrentContext: selected.ContextName}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)

	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("parse kubeconfig source %q: %w", source.Name, err)
	}

	if _, ok := rawConfig.Contexts[selected.ContextName]; !ok {
		return nil, fmt.Errorf("context %q missing from source %q", selected.ContextName, source.Name)
	}

	ensureContextNamespace(&rawConfig, selected.ContextName, selected.Namespace)
	clientConfig = clientcmd.NewNonInteractiveClientConfig(rawConfig, selected.ContextName, overrides, loadingRules)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("build rest config for %q: %w", selected.ContextName, err)
	}

	restConfig = rest.CopyConfig(restConfig)
	restConfig.Timeout = 15 * time.Second

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %w", err)
	}

	apiExtensionsClient, err := apiextensionsclient.NewForConfig(rest.CopyConfig(restConfig))
	if err != nil {
		return nil, fmt.Errorf("create api extensions client: %w", err)
	}

	return &ClusterClient{
		Context:             selected,
		Clientset:           clientset,
		APIExtensionsClient: apiExtensionsClient,
	}, nil
}

func ensureContextNamespace(config *clientcmdapi.Config, contextName, namespace string) {
	if config == nil || namespace == "" {
		return
	}

	if configContext, ok := config.Contexts[contextName]; ok {
		configContext.Namespace = namespace
	}
}
