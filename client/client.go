package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/AmirMirzayi/relay/config"
	"github.com/AmirMirzayi/relay/utils"
)

func Receive(ip net.IP, port, progressbarWidth int, timeout time.Duration) error {
	serverAddress := fmt.Sprintf("%s:%d", ip, port)
	fmt.Printf("Connecting to %s...", serverAddress)
	conn, err := net.DialTimeout("tcp", serverAddress, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Printf("\rSuccessfully Connected to %s ✓\n", serverAddress)

	files, err := receiveFileInfo(conn)
	if err != nil {
		return err
	}

	fmt.Printf("Preparing to receive %d files with %s\n", len(files), files.HumanReadableTotalSize())

	for _, file := range files {
		if err = receiveFile(conn, file, progressbarWidth); err != nil {
			return err
		}
	}

	return nil
}

func receiveFile(conn net.Conn, file config.File, progressbarWidth int) error {
	fileDir := ""
	if len(file.Parents) > 0 {
		fileDir := filepath.Join(file.Parents...)
		dir := filepath.Dir(filepath.Join(fileDir, file.Name))
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Create(filepath.Join(fileDir, file.Name))
	if err != nil {
		return fmt.Errorf("failed to create file %s, %v", file.Name, err)
	}
	defer f.Close()

	shortedFileName := utils.ShortedString(file.Name, 10, 6, 4)

	buffer := make([]byte, config.DefaultChunkSize)
	var totalBytesRead, totalBytesWritten int64
	byteToRead := config.DefaultChunkSize

	for totalBytesRead < file.Size {
		if config.DefaultChunkSize > file.Size-totalBytesRead {
			byteToRead = int(file.Size - totalBytesRead)
		}
		n, err := io.ReadFull(conn, buffer[:byteToRead])
		if err != nil {
			return fmt.Errorf("failed to read buffer from network, %v", err)
		}

		m, err := f.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf("failed to write buffer into file, %v", err)
		}

		totalBytesRead += int64(n)
		totalBytesWritten += int64(m)

		receivedPercent := int(totalBytesWritten * 100 / file.Size)
		utils.DrawProgressBar(receivedPercent, progressbarWidth, shortedFileName)
	}
	fmt.Printf("\r%s ✓\n", file.Name)

	return nil
}

func receiveFileInfo(conn net.Conn) (config.Files, error) {
	var dataLen uint32
	if err := binary.Read(conn, binary.BigEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("failed to read files info, %v", err)
	}

	buffer := make([]byte, dataLen)
	if _, err := conn.Read(buffer); err != nil {
		return nil, fmt.Errorf("failed to read files info, %v", err)
	}

	var files config.Files
	if err := json.Unmarshal(buffer, &files); err != nil {
		return nil, fmt.Errorf("failed to unmarshal files info, %v", err)
	}

	return files, nil
}
