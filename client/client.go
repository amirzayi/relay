// Package client enables file reception from a server over a network.
// It handles connecting to the host, receiving files and displaying progress.
package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/AmirMirzayi/relay/config"
	"github.com/AmirMirzayi/relay/utils"
)

// Receive connects to a file transfer server at the specified IP and port to receive files.
// It handles file reception and displays a progress bar.
func Receive(setting config.Setting) error {
	serverAddress := fmt.Sprintf("%s:%d", setting.IP, setting.Port)
	fmt.Printf("Connecting to %s...", serverAddress)
	conn, err := net.DialTimeout("tcp", serverAddress, setting.Timeout)
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

	for i, file := range files {
		if err = receiveFile(conn, file, i+1, setting); err != nil {
			return err
		}
	}

	return nil
}

func receiveFile(conn net.Conn, file config.File, fileID int, setting config.Setting) error {
	filePath := filepath.Join(file.Parents...)
	filePath = filepath.Join(setting.SavePath, filePath, file.Name)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s, %v", file.Name, err)
	}
	defer f.Close()

	if setting.SilentTransfer {
		return utils.WriteFromReader(conn, f, file.Size, setting.BufferSize)
	}

	shortedFileName := utils.ShortedString(file.Name, 10, 8, 3)
	fileSize := utils.ConvertByteSizeToHumanReadable(float64(file.Size))
	barTitle := fmt.Sprintf("<%s ^ %s>", fileSize, shortedFileName)
	if err = utils.DrawRWProgressbar(conn, f, file.Size, setting.BufferSize, setting.ProgressbarWidth, barTitle); err != nil {
		return err
	}
	fmt.Printf("\r[%d] %s ✓\033[K\n", fileID, file.Name)

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
