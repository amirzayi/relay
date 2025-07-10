package config

import (
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/amirzayi/relay/utils"
)

const (
	DefaultPort             = 55555
	DefaultTimeout          = 30 * time.Second
	DefaultBufferSize       = 1024 * 1024
	DefaultProgressbarWidth = 25
	DefaultSilent           = false

	DefaultGUIPort    = "1112"
	DefaultGUITimeout = time.Second
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

type File struct {
	Name    string
	Size    int64
	Path    string
	Parents []string
}

type Files []File

func (fs Files) TotalSize() (totalSize int64) {
	for _, f := range fs {
		totalSize += f.Size
	}
	return
}

func (fs Files) HumanReadableTotalSize() string {
	totalSize := fs.TotalSize()
	return utils.ConvertByteSizeToHumanReadable(float64(totalSize))
}
