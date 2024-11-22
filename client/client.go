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
)

func Receive(ip net.IP, port int, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	files, err := receiveFileInfo(conn)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err = receiveFile(conn, file); err != nil {
			return err
		}
	}

	return nil
}

func receiveFile(conn net.Conn, file config.File) error {
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

	buffer := make([]byte, config.DefaultChunkSize)
	var written int64 = 0
	byteToRead := config.DefaultChunkSize

	for written < file.Size {
		if config.DefaultChunkSize > file.Size-written {
			byteToRead = int(file.Size - written)
		}
		n, err := io.ReadFull(conn, buffer[:byteToRead])
		if err != nil {
			return fmt.Errorf("failed to read buffer from network, %v", err)
		}

		_, err = f.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf("failed to write buffer into file, %v", err)
		}

		written += int64(n)
	}

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
