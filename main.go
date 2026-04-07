package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	clusterapi "netkube/adapters/api"
	"netkube/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const authCookieName = "netkube_session"

type authConfig struct {
	Email         string
	Password      string
	SessionSecret []byte
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		log.Printf("no .env file loaded, using process environment")
	}

	if err := storage.EnsureConfigDir(); err != nil {
		log.Fatalf("failed to create config directory: %v", err)
	}

	auth, err := loadAuthConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(authRequired(auth))

	router.LoadHTMLFiles(
		"views/index.tmpl",
		"views/login.tmpl",
		"views/404.tmpl",
		"views/sidebar.tmpl",
		"views/clusters/overview.tmpl",
		"views/clusters/node.tmpl",
		"views/clusters/details/node.tmpl",
		"views/clusters/event.tmpl",
		"views/clusters/namespaces.tmpl",
		"views/clusters/leases.tmpl",
		"views/networking/services.tmpl",
		"views/workloads/overview.tmpl",
		"views/workloads/deployments.tmpl",
		"views/workloads/pods.tmpl",
		"views/workloads/manage/deployment.tmpl",
		"views/workloads/manage/pod.tmpl",
	)

	router.Static("/public", "./public")
	router.Static("/reference", "./reference")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", gin.H{
			"path": c.Request.URL.Path,
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	router.GET("/login", func(c *gin.Context) {
		if isAuthenticated(c, auth) {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})

	router.POST("/login", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email" form:"email"`
			Password string `json:"password" form:"password"`
		}

		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login request."})
			return
		}

		email := strings.TrimSpace(request.Email)
		password := request.Password

		if email != auth.Email || password != auth.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password."})
			return
		}

		setAuthCookie(c, auth)
		c.JSON(http.StatusOK, gin.H{"redirect": "/"})
	})

	router.POST("/logout", func(c *gin.Context) {
		clearAuthCookie(c)
		c.JSON(http.StatusOK, gin.H{"redirect": "/login"})
	})

	router.GET("/clusters/overview", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/overview", gin.H{})
	})

	router.GET("/clusters/node", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/node", gin.H{})
	})

	router.GET("/clusters/nodes", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/node", gin.H{})
	})

	router.GET("/clusters/details/node", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/details/node", gin.H{})
	})

	router.GET("/clusters/event", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/event", gin.H{})
	})

	router.GET("/clusters/events", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/event", gin.H{})
	})

	router.GET("/clusters/namespaces", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/namespaces", gin.H{})
	})

	router.GET("/clusters/leases", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/leases", gin.H{})
	})

	router.GET("/networking/services", func(c *gin.Context) {
		c.HTML(http.StatusOK, "networking/services", gin.H{})
	})

	router.GET("/workloads/overview", func(c *gin.Context) {
		c.HTML(http.StatusOK, "workloads/overview", gin.H{})
	})

	router.GET("/workloads/deployments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "workloads/deployments", gin.H{})
	})

	router.GET("/workloads/pods", func(c *gin.Context) {
		c.HTML(http.StatusOK, "workloads/pods", gin.H{})
	})

	router.GET("/workloads/manage/pod", func(c *gin.Context) {
		c.HTML(http.StatusOK, "workloads/manage/pod", gin.H{})
	})

	router.GET("/workloads/manage/deployment", func(c *gin.Context) {
		c.HTML(http.StatusOK, "workloads/manage/deployment", gin.H{})
	})

	api := router.Group("/api")
	{
		api.GET("/config/sources", storage.GetSources)
		api.POST("/config/sources", storage.SaveSources)

		api.GET("/config/selected-contexts", storage.GetSelectedContexts)
		api.POST("/config/selected-contexts", storage.SaveSelectedContexts)

		api.GET("/config/contexts", storage.GetContexts)

		api.GET("/cluster/overview", clusterapi.ClusterOverviewHandler)
		api.GET("/cluster/nodes", clusterapi.ClusterNodesHandler)
		api.GET("/cluster/node", clusterapi.ClusterNodeDetailHandler)
		api.GET("/cluster/event", clusterapi.ClusterEventDetailHandler)
		api.GET("/cluster/namespaces", clusterapi.ClusterNamespacesHandler)
		api.GET("/cluster/leases", clusterapi.ClusterLeasesHandler)
		api.GET("/cluster/namespace/yaml", clusterapi.ClusterNamespaceYAMLHandler)
		api.GET("/cluster/lease/yaml", clusterapi.ClusterLeaseYAMLHandler)
		api.GET("/networking/services", clusterapi.NetworkingServicesHandler)
		api.GET("/workloads/overview", clusterapi.WorkloadsOverviewHandler)
		api.GET("/workloads/pods", clusterapi.PodsHandler)
		api.GET("/workloads/pod", clusterapi.PodDetailHandler)
		api.GET("/workloads/pod/logs", clusterapi.PodLogsHandler)
		api.GET("/workloads/pod/events", clusterapi.PodEventsHandler)
		api.GET("/workloads/pod/yaml", clusterapi.PodYAMLHandler)
		api.GET("/workloads/deployments", clusterapi.DeploymentsHandler)
		api.GET("/workloads/deployment", clusterapi.DeploymentDetailHandler)
		api.GET("/workloads/deployment/events", clusterapi.DeploymentEventsHandler)
		api.GET("/workloads/deployment/yaml", clusterapi.DeploymentYAMLHandler)
	}

	log.Println("Server running in release mode on http://localhost:3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

func loadAuthConfig() (authConfig, error) {
	rawEnv, err := readRawDotEnv(".env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return authConfig{}, err
	}

	email := firstNonEmpty(rawEnv["EMAIL"], rawEnv["USERNAME"], strings.TrimSpace(os.Getenv("EMAIL")), strings.TrimSpace(os.Getenv("USERNAME")))
	password := firstNonEmpty(rawEnv["PASSWORD"], os.Getenv("PASSWORD"))
	secret := firstNonEmpty(rawEnv["SESSION_SECRET"], os.Getenv("SESSION_SECRET"))

	if email == "" {
		email = strings.TrimSpace(os.Getenv("USERNAME"))
	}

	if email == "" || password == "" {
		return authConfig{}, errors.New("missing auth credentials: set EMAIL (or USERNAME) and PASSWORD in the environment")
	}

	if secret == "" {
		sum := sha256.Sum256([]byte(email + "\x00" + password))
		secret = base64.StdEncoding.EncodeToString(sum[:])
	}

	return authConfig{
		Email:         email,
		Password:      password,
		SessionSecret: []byte(secret),
	}, nil
}

func authRequired(auth authConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if isPublicPath(path) {
			c.Next()
			return
		}

		if isAuthenticated(c, auth) {
			c.Next()
			return
		}

		if strings.HasPrefix(path, "/api/") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required."})
			return
		}

		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
	}
}

func isPublicPath(path string) bool {
	return path == "/login" ||
		strings.HasPrefix(path, "/public/") ||
		strings.HasPrefix(path, "/reference/")
}

func isAuthenticated(c *gin.Context, auth authConfig) bool {
	cookie, err := c.Cookie(authCookieName)
	if err != nil {
		return false
	}

	parts := strings.Split(cookie, ".")
	if len(parts) != 2 {
		return false
	}

	emailBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	email := string(emailBytes)
	if email != auth.Email {
		return false
	}

	expected := signSessionValue(email, auth.SessionSecret)
	return hmac.Equal(signature, expected)
}

func setAuthCookie(c *gin.Context, auth authConfig) {
	email := base64.RawURLEncoding.EncodeToString([]byte(auth.Email))
	signature := base64.RawURLEncoding.EncodeToString(signSessionValue(auth.Email, auth.SessionSecret))
	value := email + "." + signature

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(authCookieName, value, 60*60*24*7, "/", "", false, true)
}

func clearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(authCookieName, "", -1, "/", "", false, true)
}

func signSessionValue(value string, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(value))
	return h.Sum(nil)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

func readRawDotEnv(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}

		values[key] = trimEnvQuotes(value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func trimEnvQuotes(value string) string {
	if len(value) >= 2 {
		if (value[0] == '\'' && value[len(value)-1] == '\'') || (value[0] == '"' && value[len(value)-1] == '"') {
			return value[1 : len(value)-1]
		}
	}

	return value
}
