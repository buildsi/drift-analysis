package main

import (
	"github.com/buildsi/drift-server/config"
	"github.com/buildsi/drift-server/models"
	"github.com/buildsi/drift-server/urls"
)

// @title Drift Server API
// @version 1.0
// @description Record metadata about inflection point changes in packages
// @termsOfService http://swagger.io/terms/

// @contact.name @vsoch
// @contact.url https://github.com/buildsi/drift-server/issues

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @securityDefinitions.basic BasicAuth
func main() {

	// Create a new config to get envars
	c := config.NewConfig()

	// Connect and migrate the database
	models.InitDatabase(&c)

	// Prepare urls and authentication
	urls.Prepare(&c)
}
