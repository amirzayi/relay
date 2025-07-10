// Package host provides functionality for serving files to clients over a network.
// It manage file transfer sessions, network configurations and client interactions.
package host

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/amirzayi/relay/config"
	"github.com/amirzayi/relay/pkg/fileutil"
	"github.com/amirzayi/relay/utils"
)

// Serve starts a file transfer server that listens on the specified IP and port.
// It serves files located at the provided paths to connected clients.
func Serve(setting config.Setting, paths ...string) error {
	files, err := fileutil.GetFilesByPaths(paths...)
	if err != nil {
		return err
	}

	fmt.Printf("Preparing to send %d files with %s\n", len(files), files.HumanReadableTotalSize())

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: setting.IP, Port: setting.Port})
	if err != nil {
		return err
	}
	defer listener.Close()
	listener.SetDeadline(time.Now().Add(setting.Timeout))
	conn, err := listener.AcceptTCP()
	// conn, err := listenWithTimeout(listener, setting.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = sendFilesInfo(conn, files); err != nil {
		return err
	}

	for i, file := range files {
		if err = sendFile(conn, file, i+1, setting); err != nil {
			return err
		}
	}

	return nil
}

func sendFilesInfo(conn net.Conn, files config.Files) error {
	return gob.NewEncoder(conn).Encode(files)
}

func sendFile(conn net.Conn, file config.File, fileID int, setting config.Setting) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("failed to load file %s, %v", file.Path, err)
	}
	defer f.Close()

	if setting.SilentTransfer {
		return utils.WriteFromReader(f, conn, file.Size, setting.BufferSize)
	}

	shortedFileName := utils.ShortedString(file.Name, 10, 8, 3)
	fileSize := utils.ConvertByteSizeToHumanReadable(float64(file.Size))
	barTitle := fmt.Sprintf("<%s ^ %s>", fileSize, shortedFileName)
	if err = utils.DrawRWProgressbar(f, conn, file.Size, setting.BufferSize, setting.ProgressbarWidth, barTitle); err != nil {
		return err
	}
	fmt.Printf("\r[%d] %s ✓\033[K\n", fileID, file.Name)

	return nil
}

func listenWithTimeout(listener net.Listener, dur time.Duration) (net.Conn, error) {
	done := make(chan error)

	var (
		conn net.Conn
		err  error
	)
	go func() {
		conn, err = listener.Accept()
		done <- err
	}()

	select {

	case <-time.After(dur):
		return nil, fmt.Errorf("no connection accepted after %.1f seconds", dur.Seconds())
	case <-done:
		return conn, err
	}
}
