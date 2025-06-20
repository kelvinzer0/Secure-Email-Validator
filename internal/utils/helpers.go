package utils

import (
	"fmt"
	"runtime"
)

// GetVersion returns the application version
func GetVersion() string {
	return "1.0.0"
}

// GetBuildInfo returns build information
func GetBuildInfo() string {
	return fmt.Sprintf("Go %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// IsValidDomain performs basic domain validation
func IsValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	
	// Basic domain format check
	if domain[0] == '.' || domain[len(domain)-1] == '.' {
		return false
	}
	
	return true
}
