package urls

import (
	"fmt"

	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"

	_ "github.com/buildsi/drift-server/docs"

	"github.com/buildsi/drift-server/config"
	"github.com/buildsi/drift-server/views"
	"github.com/gin-gonic/gin"
)

func Prepare(c *config.Config) {

	// Use gin to create the server
	router := gin.Default()

	// All views are authenticated
	if (c.Username != "") && (c.Password != "") {

		fmt.Println("Server will use authentication.")

		// Authentication required group using gin.BasicAuth() middleware
		authorized := router.Group("/", gin.BasicAuth(gin.Accounts{c.Username: c.Password}))
		authorized.GET("/", views.ServerInfo)

		authorized.POST("/spec/", views.Spec)
		authorized.POST("/inflection-point/", views.InflectionPoint)
		authorized.POST("/build/", views.Build)

		authorized.GET("/inflection-points/", views.ListInflectionPoints)
		authorized.GET("/inflection-points/:package", views.ListInflectionPoints)
		authorized.GET("/inflection-points/:package/*version", views.ListInflectionPoints)

		authorized.GET("/commits/", views.ListCommits)
		authorized.GET("/packages/", views.ListPackages)
		authorized.GET("/tags/", views.ListTags)
		authorized.GET("/specs/", views.ListSpecs)
		authorized.GET("/builds/", views.ListBuilds)
		authorized.GET("/api/", views.Swagger)
		authorized.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	} else {
		router.GET("/", views.ServerInfo)

		router.POST("/spec/", views.Spec)
		router.POST("/inflection-point/", views.InflectionPoint)
		router.POST("/build/", views.Build)

		router.GET("/inflection-points/", views.ListInflectionPoints)
		router.GET("/inflection-points/:package", views.ListInflectionPoints)
		router.GET("/inflection-points/:package/*version", views.ListInflectionPoints)

		router.GET("/commits/", views.ListCommits)
		router.GET("/packages/", views.ListPackages)
		router.GET("/tags/", views.ListTags)
		router.GET("/specs/", views.ListSpecs)
		router.GET("/builds/", views.ListBuilds)
		router.GET("/api/", views.Swagger)
		router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	}

	router.Run(":" + c.Port)

}
