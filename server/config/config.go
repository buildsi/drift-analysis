package config

import (
	"os"
	"reflect"
	"strings"
)

type Config struct {
	API api
	DB  db
	S3  s3
}

type api struct {
	Username string
	Password string
	Port     string
}

type db struct {
	Path string
}

type s3 struct {
	Bucket    string
	AccessKey string
	SecretKey string
	Endpoint  string
	Region    string
}

// Init creates a new config based on parsed environment variables.
func Init() (*Config, error) {
	return envParseConfig(&Config{
		API: api{
			Username: "",
			Password: "",
			Port:     "8080",
		},
		DB: db{
			Path: "points.db",
		},
		S3: s3{
			Bucket:    "",
			AccessKey: "",
			SecretKey: "",
			Endpoint:  "",
			Region:    "us-east-1",
		},
	}), nil
}

func envParseConfig(in *Config) *Config {
	numSubStructs := reflect.ValueOf(in).Elem().NumField()
	for i := 0; i < numSubStructs; i++ {
		iter := reflect.ValueOf(in).Elem().Field(i)
		subStruct := strings.ToUpper(iter.Type().Name())

		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			fieldName := structType.Field(j).Name
			for _, prefix := range []string{"DRIFTSERVER"} {
				evName := prefix + "_" + subStruct + "_" + strings.ToUpper(fieldName)
				evVal, evExists := os.LookupEnv(evName)
				if evExists && evVal != fieldVal {
					iter.FieldByName(fieldName).SetString(evVal)
				}
			}
		}
	}
	return in
}
