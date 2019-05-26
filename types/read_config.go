// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// Package types contains definitions for public types
package types

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// OsEnv implements interface to wrap os.Getenv
type OsEnv struct {
}

// Getenv wraps os.Getenv
func (OsEnv) Getenv(key string) string {
	return os.Getenv(key)
}

// HasEnv provides interface for os.Getenv
type HasEnv interface {
	Getenv(key string) string
}

// ReadConfig constitutes config from env variables
type ReadConfig struct {
}

func parseIntValue(val string, fallback int) int {
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return parsedVal
		}
	}
	return fallback
}

func parseIntOrDurationValue(val string, fallback time.Duration) time.Duration {
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return time.Duration(parsedVal) * time.Second
		}
	}

	duration, durationErr := time.ParseDuration(val)
	if durationErr != nil {
		return fallback
	}

	return duration
}

func parseBoolValue(val string, fallback bool) bool {
	if len(val) > 0 {
		return val == "true"
	}
	return fallback
}

func parseString(val string, fallback string) string {
	if len(val) > 0 {
		return val
	}
	return fallback
}

// Read fetches config from environmental variables.
func (r ReadConfig) Read(hasEnv HasEnv) BootstrapConfig {
	r.ensureRequired()
	defaultTCPPort := 8080

	cfg := BootstrapConfig{}

	cfg.ReadTimeout = parseIntOrDurationValue(hasEnv.Getenv("read_timeout"), time.Minute*3)
	cfg.WriteTimeout = parseIntOrDurationValue(hasEnv.Getenv("write_timeout"), time.Minute*3)
	cfg.Port = parseIntValue(hasEnv.Getenv("port"), defaultTCPPort)

	providers := strings.Split(os.Getenv("providers"), ",")
	for i, v := range providers {
		providers[i] = strings.Trim(v, " ")
	}

	cfg.Providers = providers
	cfg.DefaultProvider = os.Getenv("default_provider")
	return cfg
}

func (r ReadConfig) ensureRequired() {
	var required = []string{"providers", "default_provider"}

	for _, v := range required {
		if _, ok := os.LookupEnv(v); !ok {
			panic(fmt.Sprintf("%s environment variable must be set, please see README.md", v))
		}
	}
}

// BootstrapConfig for the process.
type BootstrapConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	Providers       []string
	DefaultProvider string
}
