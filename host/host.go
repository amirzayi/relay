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
	"github.com/AmirMirzayi/relay/utils"
)

func Serve(ip net.IP, port, progressbarWidth int, timeout time.Duration, pathes ...string) error {
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
		if err = sendFile(conn, file, progressbarWidth); err != nil {
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
				Name:    fileInfo.Name(),
				Size:    fileInfo.Size(),
				Path:    path,
				Parents: nil,
			})
			continue
		}

		dirFiles, err := readDirectoryFilesRecursively(path)
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

func sendFile(conn net.Conn, file config.File, progressbarWidth int) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("failed to load file %s, %v", file.Path, err)
	}
	defer f.Close()

	shortedFileName := utils.ShortedString(file.Name, 10, 6, 4)

	var totalByteSent int64
	buffer := make([]byte, config.DefaultChunkSize)
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("failed to buffering file %s, %v", file.Path, err)
		}

		byteSent, err := conn.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf("failed to write over network, %v", err)
		}

		totalByteSent += int64(byteSent)
		sentPercent := int(totalByteSent * 100 / file.Size)

		utils.DrawProgressBar(sentPercent, progressbarWidth, shortedFileName)
	}
	fmt.Printf("\r%s ✓\n", file.Name)

	return nil
}

// readDirectoryFilesRecursively retrieve directory and subdirectories files
func readDirectoryFilesRecursively(path string, parents ...string) (config.Files, error) {
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
			files = append(files, config.File{
				Name:    entry.Name(),
				Size:    entryInfo.Size(),
				Path:    filepath.Join(path, entry.Name()),
				Parents: append(parents, filepath.Base(path)),
			})
			continue
		}

		parents = append(parents, filepath.Base(path))
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
