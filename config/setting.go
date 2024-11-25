package config

import (
	"net"
	"time"
)

type Setting struct {
	IP net.IP
	Port,
	ProgressbarWidth,
	BufferSize int
	Timeout        time.Duration
	SavePath       string
	SilentTransfer bool
}

func DefaultSetting() Setting {
	return Setting{
		IP:               net.ParseIP("127.0.0.1"),
		Port:             DefaultPort,
		ProgressbarWidth: DefaultProgressbarWidth,
		Timeout:          DefaultTimeout,
		SavePath:         DefaultDirectory(),
		SilentTransfer:   DefaultSilent,
	}
}
