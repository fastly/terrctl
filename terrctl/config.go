package main

import (
	"flag"
	"time"
)

// ConfigStruct - Structure to store the configuration
type ConfigStruct struct {
	DeployTimeout       time.Duration
	HealthTimeout       time.Duration
	HTTPClientTimeout   time.Duration
	Language            string
	MaxDeployAttempts   uint
	MaxResponseBodySize int64
}

var config = ConfigStruct{
	DeployTimeout:       90 * time.Second,
	HealthTimeout:       30 * time.Second,
	HTTPClientTimeout:   30 * time.Second,
	Language:            "auto",
	MaxDeployAttempts:   10,
	MaxResponseBodySize: 4096,
}

// Config - Return the global configuration
func Config() *ConfigStruct {
	return &config
}

// UpdateConfigFromFlags - Parse the command-line flags and update the configuration
func UpdateConfigFromFlags() error {
	language := flag.String("language", config.Language, "language (auto|c|rust|assemblyscript|wasm)")
	deployTimeout := flag.Uint64("deploy-timeout", uint64(config.DeployTimeout.Seconds()), "Timeout for deployment (seconds)")
	healthTimeout := flag.Uint64("health-timeout", uint64(config.HealthTimeout.Seconds()), "Timeout for health checks (seconds)")
	httpClientTimeout := flag.Uint64("http-timeout", uint64(config.HTTPClientTimeout.Seconds()), "Timeout for HTTP client queries (seconds)")
	maxDeployAttempts := flag.Uint("max-deploy-attempts", config.MaxDeployAttempts, "Maximum number of attempts for deployment")
	flag.Parse()
	config.Language = *language
	config.DeployTimeout = time.Duration(*deployTimeout) * time.Second
	config.HealthTimeout = time.Duration(*healthTimeout) * time.Second
	config.HTTPClientTimeout = time.Duration(*httpClientTimeout) * time.Second
	config.MaxDeployAttempts = *maxDeployAttempts
	return nil
}
