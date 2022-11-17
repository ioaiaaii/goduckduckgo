package config

import (
	"time"
)

// Init default values
const (
	DefaultQueryHTTPPort                  = "8080"
	DefaultQueryEnv                       = "production"
	DefaultQueryHTTPServerTimeout         = 30
	DefaultQueryHTTPServerShutdownTimeout = 5
	DefaultStoreDBPort                    = "5432"
	DefaultStoreDBHost                    = "localhost"
	DefaultStoreDBUser                    = ""
	DefaultStoreDBPassword                = ""
	DefaultStoreDBName                    = ""
)

type queryConfig struct {
	HTTPPort                  string
	HTTPServerTimeout         time.Duration
	HTTPServerShutdownTimeout time.Duration
}

type storeConfig struct {
	DBPort     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
}

type Config struct {
	QueryConfig queryConfig
	StoreConfig storeConfig
}
