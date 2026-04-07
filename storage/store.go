package storage

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const (
	configDir            = "config"
	uploadedSourcesDir   = "config/uploaded-sources"
	sourcesFile          = "config/sources.json"
	selectedContextsFile = "config/selected-contexts.json"
)

func EnsureConfigDir() error {
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	if err := os.MkdirAll(uploadedSourcesDir, 0o755); err != nil {
		return err
	}

	if err := ensureDefaultSourcesFile(); err != nil {
		return err
	}

	if err := migrateStoredSourcesFile(); err != nil {
		return err
	}

	return nil
}

func GetSources(c *gin.Context) {
	var payload SourcesPayload

	err := readJSONFile(sourcesFile, &payload)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusOK, SourcesPayload{Sources: []SourceItem{}})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read sources"})
		return
	}

	if payload.Sources == nil {
		payload.Sources = []SourceItem{}
	}

	c.JSON(http.StatusOK, payload)
}

func GetContexts(c *gin.Context) {
	contexts, err := LoadContexts()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusOK, ContextsPayload{Contexts: []ContextItem{}})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read sources"})
		return
	}

	c.JSON(http.StatusOK, ContextsPayload{Contexts: contexts})
}

func LoadSources() ([]SourceItem, error) {
	var payload SourcesPayload

	err := readJSONFile(sourcesFile, &payload)
	if err != nil {
		return nil, err
	}

	payload = sanitizeSourcesPayload(payload)
	if payload.Sources == nil {
		return []SourceItem{}, nil
	}

	return payload.Sources, nil
}

func LoadContexts() ([]ContextItem, error) {
	sources, err := LoadSources()
	if err != nil {
		return nil, err
	}

	contexts := buildContextsFromSources(sources)
	if contexts == nil {
		return []ContextItem{}, nil
	}

	return contexts, nil
}

func LoadStoredSourceContent(source SourceItem) ([]byte, error) {
	fullPath, err := StoredSourcePath(source)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(filepath.Clean(fullPath))
}

func StoredSourcePath(source SourceItem) (string, error) {
	if source.StoredFile == "" {
		return "", fmt.Errorf("source %s has no stored file", source.Name)
	}

	return filepath.Join(uploadedSourcesDir, source.StoredFile), nil
}

func ensureDefaultSourcesFile() error {
	_, err := os.Stat(filepath.Clean(sourcesFile))
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return writeJSONFile(sourcesFile, SourcesPayload{Sources: []SourceItem{}})
}

func migrateStoredSourcesFile() error {
	var payload SourcesPayload
	err := readJSONFile(sourcesFile, &payload)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	payload = sanitizeSourcesPayload(payload)
	return writeJSONFile(sourcesFile, payload)
}

func SaveSources(c *gin.Context) {
	var payload SourcesPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sources payload"})
		return
	}

	if payload.Sources == nil {
		payload.Sources = []SourceItem{}
	}

	payload = sanitizeSourcesPayload(payload)

	if err := syncStoredSourceFiles(payload.Sources); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save uploaded source files"})
		return
	}

	payload = stripRawContent(payload)

	if err := writeJSONFile(sourcesFile, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save sources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

func sanitizeSourcesPayload(payload SourcesPayload) SourcesPayload {
	cleaned := SourcesPayload{Sources: make([]SourceItem, 0, len(payload.Sources))}

	for _, source := range payload.Sources {
		source.ID = strings.TrimSpace(source.ID)
		source.Name = strings.TrimSpace(source.Name)
		source.Type = strings.TrimSpace(source.Type)

		if source.ID == "" || source.Name == "" {
			continue
		}

		if source.Type == "" {
			source.Type = "file"
		}

		source.StoredFile = storedFilenameForSource(source.ID, source.Name)
		cleaned.Sources = append(cleaned.Sources, source)
	}

	return cleaned
}

func stripRawContent(payload SourcesPayload) SourcesPayload {
	cleaned := SourcesPayload{Sources: make([]SourceItem, 0, len(payload.Sources))}

	for _, source := range payload.Sources {
		source.Content = ""
		cleaned.Sources = append(cleaned.Sources, source)
	}

	return cleaned
}

func syncStoredSourceFiles(sources []SourceItem) error {
	expected := make(map[string]struct{}, len(sources))

	for _, source := range sources {
		if source.StoredFile == "" {
			continue
		}

		expected[source.StoredFile] = struct{}{}
		fullPath := filepath.Join(uploadedSourcesDir, source.StoredFile)

		if strings.TrimSpace(source.Content) != "" {
			if err := os.WriteFile(filepath.Clean(fullPath), []byte(source.Content), 0o600); err != nil {
				return err
			}
			continue
		}

		if _, err := os.Stat(fullPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("missing stored file for source %s", source.Name)
			}
			return err
		}
	}

	entries, err := os.ReadDir(uploadedSourcesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if _, ok := expected[entry.Name()]; ok {
			continue
		}

		if err := os.Remove(filepath.Join(uploadedSourcesDir, entry.Name())); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	return nil
}

func storedFilenameForSource(id, name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		ext = ".yaml"
	}

	safeID := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-' || r == '_':
			return r
		default:
			return '_'
		}
	}, id)

	return safeID + ext
}

func GetSelectedContexts(c *gin.Context) {
	var payload SelectedContextsPayload

	err := readJSONFile(selectedContextsFile, &payload)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusOK, SelectedContextsPayload{SelectedContextIDs: []string{}})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read selected contexts"})
		return
	}

	if payload.SelectedContextIDs == nil {
		payload.SelectedContextIDs = []string{}
	}

	c.JSON(http.StatusOK, payload)
}

func SaveSelectedContexts(c *gin.Context) {
	var payload SelectedContextsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid selected contexts payload"})
		return
	}

	if payload.SelectedContextIDs == nil {
		payload.SelectedContextIDs = []string{}
	}

	if err := writeJSONFile(selectedContextsFile, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save selected contexts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

func buildContextsFromSources(sourceItems []SourceItem) []ContextItem {
	results := make([]ContextItem, 0)

	for _, source := range sourceItems {
		if source.StoredFile == "" {
			continue
		}

		fullPath := filepath.Join(uploadedSourcesDir, source.StoredFile)
		data, err := os.ReadFile(filepath.Clean(fullPath))
		if err != nil {
			continue
		}

		var doc kubeconfigDocument
		if err := yaml.Unmarshal(data, &doc); err != nil {
			continue
		}

		for _, entry := range doc.Contexts {
			contextName := entry.Name
			clusterName := entry.Context.Cluster
			userName := entry.Context.User
			if userName == "" {
				userName = contextName
			}
			namespace := entry.Context.Namespace
			if namespace == "" {
				namespace = "default"
			}

			server := ""
			for _, cluster := range doc.Clusters {
				if cluster.Name == clusterName {
					server = cluster.Cluster.Server
					break
				}
			}

			hasUser := false
			for _, user := range doc.Users {
				if user.Name == userName {
					hasUser = true
					break
				}
			}

			results = append(results, ContextItem{
				ID:          stableContextID(source.ID, contextName, clusterName, userName, namespace),
				ContextName: contextName,
				ClusterName: clusterName,
				UserName:    userName,
				Namespace:   namespace,
				SourceID:    source.ID,
				SourceName:  source.Name,
				Server:      server,
				IsCurrent:   doc.CurrentContext == contextName,
				HasUser:     hasUser,
			})
		}
	}

	return results
}

func stableContextID(parts ...string) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%q", parts)))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

type kubeconfigDocument struct {
	CurrentContext string                 `yaml:"current-context"`
	Contexts       []kubeconfigContextRef `yaml:"contexts"`
	Clusters       []kubeconfigClusterRef `yaml:"clusters"`
	Users          []kubeconfigUserRef    `yaml:"users"`
}

type kubeconfigContextRef struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster   string `yaml:"cluster"`
		User      string `yaml:"user"`
		Namespace string `yaml:"namespace"`
	} `yaml:"context"`
}

type kubeconfigClusterRef struct {
	Name    string `yaml:"name"`
	Cluster struct {
		Server string `yaml:"server"`
	} `yaml:"cluster"`
}

type kubeconfigUserRef struct {
	Name string `yaml:"name"`
}

func readJSONFile(path string, target any) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Clean(path), data, 0o600)
}
