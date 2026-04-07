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
		"views/sidebar.tmpl",
		"views/clusters/overview.tmpl",
		"views/workloads/overview.tmpl",
		"views/workloads/deployments.tmpl",
		"views/workloads/pods.tmpl",
	)

	router.Static("/public", "./public")
	router.Static("/reference", "./reference")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	router.GET("/clusters/overview", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clusters/overview", gin.H{})
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

	api := router.Group("/api")
	{
		api.GET("/config/sources", storage.GetSources)
		api.POST("/config/sources", storage.SaveSources)

		api.GET("/config/selected-contexts", storage.GetSelectedContexts)
		api.POST("/config/selected-contexts", storage.SaveSelectedContexts)

		api.GET("/config/contexts", storage.GetContexts)

		api.GET("/cluster/overview", clusterapi.ClusterOverviewHandler)
		api.GET("/workloads/overview", clusterapi.WorkloadsOverviewHandler)
		api.GET("/workloads/pods", clusterapi.PodsHandler)
		api.GET("/workloads/deployments", clusterapi.DeploymentsHandler)
	}

	log.Println("Server running in release mode on http://localhost:3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
