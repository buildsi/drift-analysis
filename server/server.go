package main

import (
	"fmt"
	"log"

	"github.com/buildsi/drift-analysis/server/config"
	"github.com/buildsi/drift-analysis/server/database"
	"github.com/buildsi/drift-analysis/server/datastore"
	"github.com/buildsi/drift-analysis/server/web"
)

func main() {
	var auth map[string]string

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Init(cfg.DB.Path)
	if err != nil {
		log.Fatal(err)
	}

	s3Opts := datastore.S3ConfigOptions{
		Bucket:    cfg.S3.Bucket,
		Endpoint:  cfg.S3.Endpoint,
		Region:    cfg.S3.Region,
		AccessKey: cfg.S3.AccessKey,
		SecretKey: cfg.S3.SecretKey,
	}

	ds, err := datastore.Init(s3Opts)

	// check if authentication creds are defined
	if cfg.API.Username != "" {
		auth = make(map[string]string)
		auth[cfg.API.Username] = cfg.API.Password
	}

	err = web.Start(fmt.Sprintf(":%s", cfg.API.Port), auth, db, ds)
	if err != nil {
		log.Fatal(err)
	}
}
