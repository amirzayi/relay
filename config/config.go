package config

import (
	"net"
	"time"
)

const (
	DefaultPort      = 55555
	DefaultTimeout   = 30 * time.Second
	DefaultChunkSize = 1024 * 1024
	ProgressbarWidth = 25
)

func DefaultIP() net.IP {
	return net.IPv4zero
}

type File struct {
	Name    string   `json:"name"`
	Size    int64    `json:"size"`
	Path    string   `json:"-"`
	Parents []string `json:"parents"`
}

type Files []File
