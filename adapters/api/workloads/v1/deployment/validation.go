package deployment

import (
	"strings"

	"netkube/adapters/api/shared"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
)

func decodeManifest(content string) (*appsv1.Deployment, error) {
	manifest := strings.TrimSpace(content)
	if manifest == "" {
		return nil, shared.ErrManifestRequired("deployment")
	}

	var deployment appsv1.Deployment
	if err := yaml.Unmarshal([]byte(manifest), &deployment); err != nil {
		return nil, shared.ErrManifestInvalid("deployment", err)
	}

	if !strings.EqualFold(deployment.Kind, "Deployment") {
		return nil, shared.ErrManifestKind("Deployment")
	}
	if deployment.APIVersion != "apps/v1" {
		return nil, shared.ErrManifestAPIVersion("apps/v1")
	}
	if strings.TrimSpace(deployment.Name) == "" {
		return nil, shared.ErrManifestNameRequired("deployment")
	}
	if strings.TrimSpace(deployment.Namespace) == "" {
		deployment.Namespace = "default"
	}

	return &deployment, nil
}
