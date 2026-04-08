package pod

import (
	"strings"

	"netkube/adapters/api/shared"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func decodeManifest(content string) (*corev1.Pod, error) {
	manifest := strings.TrimSpace(content)
	if manifest == "" {
		return nil, shared.ErrManifestRequired("pod")
	}

	var pod corev1.Pod
	if err := yaml.Unmarshal([]byte(manifest), &pod); err != nil {
		return nil, shared.ErrManifestInvalid("pod", err)
	}

	if !strings.EqualFold(pod.Kind, "Pod") {
		return nil, shared.ErrManifestKind("Pod")
	}
	if pod.APIVersion != "v1" {
		return nil, shared.ErrManifestAPIVersion("v1")
	}
	if strings.TrimSpace(pod.Name) == "" {
		return nil, shared.ErrManifestNameRequired("pod")
	}
	if strings.TrimSpace(pod.Namespace) == "" {
		pod.Namespace = "default"
	}

	return &pod, nil
}
