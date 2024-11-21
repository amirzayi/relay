package config

import (
	"time"
)

const (
	DefaultIP        = "127.0.0.1"
	DefaultPort      = 55555
	DefaultTimeout   = 30 * time.Second
	DefaultChunkSize = 1024 * 1024
)

type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"-"`
}

type Files []File
