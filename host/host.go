package host

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/AmirMirzayi/relay/config"
)

func Serve(ip net.IP, port int, timeout time.Duration, pathes ...string) error {
	files, err := getFilesByPathes(pathes...)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: ip, Port: port})
	if err != nil {
		return err
	}
	defer listener.Close()

	conn, err := listenWithTimeout(listener, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = sendFilesInfo(conn, files); err != nil {
		return err
	}

	for _, file := range files {
		if err = sendFile(conn, file); err != nil {
			return err
		}
	}

	return nil
}

func getFilesByPathes(pathes ...string) (config.Files, error) {
	// preallocate files to args lengths
	// but, what if an arg was directory
	files := make(config.Files, 0, len(pathes))

	for _, path := range pathes {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read information of %s, %v", path, err)
		}

		if !fileInfo.IsDir() {
			files = append(files, config.File{
				Name: fileInfo.Name(),
				Size: fileInfo.Size(),
				Path: path,
			})
			continue
		}

		// read directory files recursively
		dirFiles, err := readDirectoryFiles(path)
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

func sendFile(conn net.Conn, file config.File) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("failed to load file %s, %v", file.Path, err)
	}
	defer f.Close()

	buffer := make([]byte, config.DefaultChunkSize)
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("failed to buffering file %s, %v", file.Path, err)
		}
		if _, err = conn.Write(buffer[:n]); err != nil {
			return fmt.Errorf("failed to write over network, %v", err)
		}
	}

	return nil
}

func readDirectoryFiles(path string, parent ...string) (config.Files, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve directory data on %s, %v", path, err)
	}

	var files []config.File

	for _, entry := range entries {
		entryInfo, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to load info of %s, %v", entry.Name(), err)
		}

		if !entryInfo.IsDir() {
			fPath := filepath.Join(parent...)
			fPath = filepath.Join(fPath, filepath.Base(path), entry.Name())
			files = append(files, config.File{
				Name: fPath,
				Size: entryInfo.Size(),
				Path: filepath.Join(path, entry.Name()),
			})
			continue
		}

		innerFiles, err := readDirectoryFiles(filepath.Join(path, entry.Name()), filepath.Base(path))
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
