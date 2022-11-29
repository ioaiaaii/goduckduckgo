package config

import (
	"time"
)

// Init default values
const (
	DefaultDDGEndpoint                    = "https://api.duckduckgo.com"
	DefaultQueryHTTPPort                  = "8080"
	DefaultQueryEnv                       = "production"
	DefaultQueryHTTPServerTimeout         = 30
	DefaultQueryHTTPServerShutdownTimeout = 5
	DefaultStoreDBPort                    = "5432"
	DefaultStoreDBHost                    = "localhost"
	DefaultStoreDBUser                    = ""
	DefaultStoreDBPassword                = ""
	DefaultStoreDBName                    = ""
	DefaultStoreGRPCAddress               = "127.0.0.1"
	DefaultStoreGRPCPort                  = "50052"
	DefaultQueryEndpoint                  = DefaultStoreGRPCAddress + ":" + DefaultStoreGRPCPort
)

type queryConfig struct {
	HTTPPort                  string
	HTTPServerTimeout         time.Duration
	HTTPServerShutdownTimeout time.Duration
	StoreEndpoint             string
}

type storeConfig struct {
	DBPort           string
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	StoreGRPCAddress string
	StoreGRPCPort    string
}

type Config struct {
	QueryConfig queryConfig
	StoreConfig storeConfig
}
