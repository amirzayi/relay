package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func Serve() error {
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		return err
	}
	defer listener.Close()

	conn, err := ListenWithTimeout(listener, 30*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	fileName := "Reza Yazdani - Chelcheragh.mp3"

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	fi := FileInfo{
		Name: fileInfo.Name(),
		Size: fileInfo.Size(),
	}
	bytes, err := json.Marshal(fi)
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

	if _, err = io.Copy(conn, file); err != nil {
		return err
	}

	return nil
}

func ListenWithTimeout(listener net.Listener, dur time.Duration) (net.Conn, error) {
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
