package views

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/buildsi/drift-server/models"
	"github.com/buildsi/drift-server/version"
	"github.com/gin-gonic/gin"
)

// General function to return error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// ServerInfo godoc
// @Summary Get server info
// @Description get server info
// @ID get-server-info
// @Produce  json
// @Success 200
// @Router / [get]
func ServerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": version.BuildVersion, "status": "running"})
}

// Spec godoc
// @Summary Create a new package spec
// @Description create a new package spec
// @ID post-spec
// @Accept  json
// @Produce  json
// @Success 200
// @Router /spec/ [post]
func Spec(ctx *gin.Context) {

	// The same as an inflection point, but w/ json and only core fields
	specRequest := new(models.SpecRequest)
	if err := ctx.ShouldBind(&specRequest); err != nil {
		fmt.Println("Error parsing", err)
	}
	// Get or Create the spec
	spec := models.GetOrCreateSpec(specRequest)

	// Return the unique id
	ctx.JSON(http.StatusOK, gin.H{"ID": spec.ID})
}

// Build godoc
// @Summary Create a new build
// @Description create a new build
// @ID post-build
// @Accept  json
// @Produce  json
// @Success 200
// @Router /build/ [post]
func Build(ctx *gin.Context) {

	// The same as an inflection point, but w/ json and only core fields
	buildRequest := new(models.BuildRequest)
	if err := ctx.ShouldBind(&buildRequest); err != nil {
		fmt.Println("Error parsing", err)
	}
	// Get or Create the build
	build := models.GetOrCreateBuild(buildRequest)

	// Return the unique id
	ctx.JSON(http.StatusOK, gin.H{"ID": build.ID})
}

// InflectionPoints godoc
// @Summary Create a new inflection point
// @Description created a new inflection point
// @ID post-inflection-point
// @Accept  json
// @Produce  json
// @Success 200
// @Router /inflection-point/ [post]
func InflectionPoint(ctx *gin.Context) {
	// The same as an inflection point, but w/ json and only core fields
	pointRequest := models.InflectionPointRequest{}
	if err := ctx.ShouldBind(&pointRequest); err != nil {
		fmt.Println("Error parsing", err)
	}
	// Get or Create the point
	inflectionPoint := models.GetOrCreateInflectionPoint(&pointRequest)

	// Return the unique id
	ctx.JSON(http.StatusOK, gin.H{"ID": inflectionPoint.ID})
}

func getPageNumber(ctx *gin.Context) int {
	// A request to list inflection points can have page number
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return 1
	}
	if req.PageID == 0 {
		req.PageID = 1
	}
	return req.PageID
}

func getPackageInfo(ctx *gin.Context) (name string, version string) {
	// A request to list inflection points can have page number
	name = ctx.Param("package")
	version = strings.ReplaceAll(ctx.Param("version"), "/", "")
	return name, version
}

// GET /swagger/
// Redirect to /swagger/index.html
func Swagger(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	ctx.Abort()
}

// ListInflectionPoints godoc
// @Summary List all inflection points
// @Description List all inflectoin points
// @ID get-inflection-points
// @Accept  json
// @Produce  json
// @Success 200
// @Router /inflection-points/ [get]
func ListInflectionPoints(ctx *gin.Context) {
	pageID := getPageNumber(ctx)
	PackageName, PackageVersion := getPackageInfo(ctx)
	points := models.ListInflectionPoints(pageID, PackageName, PackageVersion)
	fmt.Println(points)
	ctx.JSON(http.StatusOK, points)
}

// ListCommits godoc
// @Summary List all commits
// @Description List commits
// @ID get-commits
// @Accept  json
// @Produce  json
// @Success 200
// @Router /commits/ [get]
func ListCommits(ctx *gin.Context) {
	pageID := getPageNumber(ctx)
	commits := models.ListCommits(pageID)
	ctx.JSON(http.StatusOK, commits)
}

// ListCommits godoc
// @Summary List all tags
// @Description List tags
// @ID get-tags
// @Accept  json
// @Produce  json
// @Success 200
// @Router /tags/ [get]
func ListTags(ctx *gin.Context) {
	pageID := getPageNumber(ctx)
	tags := models.ListTags(pageID)
	ctx.JSON(http.StatusOK, tags)
}

// ListPackages godoc
// @Summary List all packages
// @Description List packages
// @ID get-packages
// @Accept  json
// @Produce  json
// @Success 200
// @Router /packages/ [get]
func ListPackages(ctx *gin.Context) {
	pageID := getPageNumber(ctx)
	packages := models.ListPackages(pageID)
	ctx.JSON(http.StatusOK, packages)
}

// ListSpecs godoc
// @Summary List all specs
// @Description List specs
// @ID get-specs
// @Accept  json
// @Produce  json
// @Success 200
// @Router /specs/ [get]
func ListSpecs(ctx *gin.Context) {
	pageID := getPageNumber(ctx)
	specs := models.ListSpecs(pageID)
	ctx.JSON(http.StatusOK, specs)
}

// ListBuilds godoc
// @Summary List all builds
// @Description List builds
// @ID get-builds
// @Accept  json
// @Produce  json
// @Success 200
// @Router /builds/ [get]
func ListBuilds(ctx *gin.Context) {
	pageID := getPageNumber(ctx)
	builds := models.ListBuilds(pageID)
	ctx.JSON(http.StatusOK, builds)
}
