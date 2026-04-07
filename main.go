package main

import (
	"log"
	"net/http"

	clusterapi "netkube/adapters/api"
	"netkube/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		log.Printf("no .env file loaded, using process environment")
	}

	if err := storage.EnsureConfigDir(); err != nil {
		log.Fatalf("failed to create config directory: %v", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.LoadHTMLFiles(
		"views/index.tmpl",
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
