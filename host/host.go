// Package host provides functionality for serving files to clients over a network.
// It manage file transfer sessions, network configurations and client interactions.
package host

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/AmirMirzayi/relay/config"
	"github.com/AmirMirzayi/relay/utils"
)

// Serve starts a file transfer server that listens on the specified IP and port.
// It serves files located at the provided paths to connected clients.
func Serve(setting config.Setting, paths ...string) error {
	files, err := getFilesByPaths(paths...)
	if err != nil {
		return err
	}

	fmt.Printf("Preparing to send %d files with %s\n", len(files), files.HumanReadableTotalSize())

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: setting.IP, Port: setting.Port})
	if err != nil {
		return err
	}
	defer listener.Close()

	conn, err := listenWithTimeout(listener, setting.Timeout)
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

func getFilesByPaths(paths ...string) (config.Files, error) {
	// preallocate files to args lengths
	// but, what if an arg was directory
	files := make(config.Files, 0, len(paths))

	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read information of %s, %v", path, err)
		}

		if !fileInfo.IsDir() {
			files = append(files, config.File{
				Name:    fileInfo.Name(),
				Size:    fileInfo.Size(),
				Path:    path,
				Parents: nil,
			})
			continue
		}

		dirFiles, err := readDirectoryFilesRecursively(path, filepath.Base(path))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve files on directory %s, %v", path, err)
		}
		files = append(files, dirFiles...)
	}

	return files, nil
}

func sendFilesInfo(conn net.Conn, files config.Files) error {
	bytes, err := json.Marshal(files)
	if err != nil {
		return fmt.Errorf("failed to prepare transmit information, %v", err)
	}

	dataLen := uint32(len(bytes))
	if err = binary.Write(conn, binary.BigEndian, dataLen); err != nil {
		return fmt.Errorf("failed to write data size over network, %v", err)
	}
	if _, err = conn.Write(bytes); err != nil {
		return fmt.Errorf("failed to write transmit information over network, %v", err)
	}

	return nil
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

// readDirectoryFilesRecursively retrieve directory and subdirectories files
func readDirectoryFilesRecursively(path string, parents ...string) (config.Files, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve directory data on %s, %v", path, err)
	}

	var files []config.File

	// we need to separate files and directories iteration over entries because of confusing
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		entryInfo, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to load info of %s, %v", entry.Name(), err)
		}

		files = append(files, config.File{
			Name:    entry.Name(),
			Size:    entryInfo.Size(),
			Path:    filepath.Join(path, entry.Name()),
			Parents: parents,
		})
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		parents = append(parents, entry.Name())
		innerFiles, err := readDirectoryFilesRecursively(filepath.Join(path, entry.Name()), parents...)
		if err != nil {
			return nil, err
		}
		files = append(files, innerFiles...)
	}

	return files, nil
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
