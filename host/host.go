package host

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

func Serve(ip net.IP, port int, timeout time.Duration, pathes ...string) error {
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

	files := make(config.Files, 0, len(pathes))

	for _, path := range pathes {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		fi := config.FileInfo{
			Name: fileInfo.Name(),
			Size: fileInfo.Size(),
			Path: path,
		}
		files = append(files, fi)
	}

	bytes, err := json.Marshal(files)
	if err != nil {
		return err
	}

	dataLen := uint32(len(bytes))
	err = binary.Write(conn, binary.BigEndian, dataLen)
	if err != nil {
		return err
	}

	if _, err = conn.Write(bytes); err != nil {
		return err
	}

	for _, file := range files {
		file, err := os.Open(file.Path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = io.Copy(conn, file); err != nil {
			// if errors.Is(err, io.EOF) {
			// 	continue
			// }
			return err
		}
	}

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
