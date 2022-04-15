package models

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/buildsi/drift-server/config"
)

// A global handle (bad idea?)
var DB *gorm.DB

// NewDatabaseConnection ...
func Connect(c *config.Config) {

	var err error
	if c.DatabaseType == "mysql" {
		fmt.Println("mysql database requested.")
		DB, err = gorm.Open(mysql.Open(c.DatabaseDSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("%v", fmt.Errorf("not possible to connect to MySQL: %w", err))
		}
	} else if c.DatabaseType == "postgres" {
		fmt.Println("postgres database requested.")
		DB, err = gorm.Open(postgres.Open(c.DatabaseDSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("%v", fmt.Errorf("not possible to connect to Postgres: %w", err))
		}
	} else if c.DatabaseType == "sqlserver" {
		fmt.Println("sqlserver database requested.")
		DB, err = gorm.Open(sqlserver.Open(c.DatabaseDSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("%v", fmt.Errorf("not possible to connect to SQLServer: %w", err))
		}
	} else {

		// If we get down here, not set or set to sqlite
		if c.DatabaseDSN == "" {
			c.DatabaseDSN = "drift-server.sqlite"
		}
		fmt.Println("sqlite database requested.")
		DB, err = gorm.Open(sqlite.Open(c.DatabaseDSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("%v", fmt.Errorf("not possible to connect to SQLite: %w", err))
		}
	}
}

func InitDatabase(c *config.Config) {

	Connect(c)

	// Migrate the schema
	DB.AutoMigrate(&Package{})
	DB.AutoMigrate(&Tag{})
	DB.AutoMigrate(&Commit{})
	DB.AutoMigrate(&InflectionPoint{})
	DB.AutoMigrate(&Spec{})
	DB.AutoMigrate(&Build{})
}
