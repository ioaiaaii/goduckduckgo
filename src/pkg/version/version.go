/*
Package version, contains build and runtime information
*/
package version

import (
	"os"
	"runtime"
	"time"
)

// Build info fields are set during a build
var (
	BuildVersion string
	BuildHash    string
	BuildDate    string
	GoVersion    = runtime.Version()
)

type GDDGVersion struct {
	BuildVersion string `json:"version"`
	BuildHash    string `json:"revision"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
}

var BuildInfo = &GDDGVersion{
	BuildVersion: BuildVersion,
	BuildHash:    BuildHash,
	BuildDate:    BuildDate,
	GoVersion:    GoVersion,
}

type GDDGRuntime struct {
	StartTime      time.Time `json:"startTime"`
	GoroutineCount int       `json:"goroutineCount"`
	GOMAXPROCS     int       `json:"GOMAXPROCS"`
	GOGC           string    `json:"GOGC"`
	GODEBUG        string    `json:"GODEBUG"`
}

var RuntimeInfo = &GDDGRuntime{
	StartTime:      time.Now(),
	GoroutineCount: runtime.NumGoroutine(),
	GOMAXPROCS:     runtime.GOMAXPROCS(0),
	GOGC:           os.Getenv("GOGC"),
	GODEBUG:        os.Getenv("GODEBUG"),
}
