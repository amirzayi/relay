package config

import (
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultPort             = 55555
	DefaultTimeout          = 30 * time.Second
	DefaultBufferSize       = 1024 * 1024
	DefaultProgressbarWidth = 25
	DefaultSilent           = false

	LookupTimeout = time.Second
	LookupPort    = 4444
)

func DefaultIP() net.IP {
	return net.IPv4zero
}

func DefaultDirectory() string {
	outputDir := "relay"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return outputDir
	}
	outputDir = filepath.Join(homeDir, outputDir)
	return outputDir
}
