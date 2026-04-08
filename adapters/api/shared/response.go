package shared

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"netkube/adapters"
)

type APIError struct {
	Error string `json:"error"`
}

type ManifestCreateRequest struct {
	Content string `json:"content"`
}

type CreatedResourceResponse struct {
	Meta      PageMeta `json:"meta"`
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Kind      string   `json:"kind"`
}

type PageMeta struct {
	ContextName string `json:"contextName"`
	ClusterName string `json:"clusterName"`
	UserName    string `json:"userName"`
	Namespace   string `json:"namespace,omitempty"`
	LastRefresh string `json:"lastRefresh"`
}

type WarningEvent struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Age       string `json:"age"`
}

func PageMetaFromCluster(cluster *adapters.ClusterClient, namespace string) PageMeta {
	return PageMeta{
		ContextName: cluster.Context.ContextName,
		ClusterName: cluster.Context.ClusterName,
		UserName:    cluster.Context.UserName,
		Namespace:   namespace,
		LastRefresh: time.Now().Format(time.RFC3339),
	}
}

func ErrManifestRequired(kind string) error {
	return fmt.Errorf("%s manifest is required", kind)
}

func ErrManifestInvalid(kind string, err error) error {
	return fmt.Errorf("invalid %s manifest: %w", kind, err)
}

func ErrManifestKind(expected string) error {
	return fmt.Errorf("manifest kind must be %s", expected)
}

func ErrManifestAPIVersion(expected string) error {
	return fmt.Errorf("manifest apiVersion must be %s", expected)
}

func ErrManifestNameRequired(kind string) error {
	return fmt.Errorf("%s metadata.name is required", kind)
}

func CreateStatusCode(err error) int {
	if err == nil {
		return http.StatusCreated
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return http.StatusBadGateway
	}

	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "already exists"):
		return http.StatusConflict
	case strings.Contains(message, "forbidden"):
		return http.StatusForbidden
	case strings.Contains(message, "unauthorized"):
		return http.StatusUnauthorized
	case strings.Contains(message, "not found"):
		return http.StatusNotFound
	case strings.Contains(message, "invalid") || strings.Contains(message, "required"):
		return http.StatusBadRequest
	default:
		return http.StatusBadGateway
	}
}
