package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/AmirMirzayi/relay/config"
)

func Receive(ip net.IP, port int, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	var dataLen uint32
	if err = binary.Read(conn, binary.BigEndian, &dataLen); err != nil {
		return err
	}

	buffer := make([]byte, dataLen)
	if _, err = conn.Read(buffer); err != nil {
		return err
	}

	var files config.Files
	if err = json.Unmarshal(buffer, &files); err != nil {
		return err
	}

	for _, file := range files {
		receiveFile(conn, file)
	}

	return nil
}

func receiveFile(conn net.Conn, file config.File) error {
	f, err := os.Create(file.Name)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.CopyN(f, conn, file.Size); err != nil {
		// if errors.Is(err, io.EOF) {
		// 	continue
		// }
		return err
	}

	return nil
}
