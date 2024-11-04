package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"os"
)

func Request() error {
	conn, err := net.Dial("tcp", "127.0.0.1:8090")
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

	var fi FileInfo
	if err = json.Unmarshal(buffer, &fi); err != nil {
		return err
	}

	f, err := os.Create(fi.Name)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, conn); err != nil {
		return err
	}

	return nil
}

type FileInfo struct {
	Name string `json:"n"`
	Size int64  `json:"s"`
}
