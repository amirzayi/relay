package gui

import (
	"encoding/gob"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amirzayi/relay/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Receiver struct {
	listener   *net.TCPListener
	connection *net.TCPConn
	files      config.Files
	chDone     chan struct{}
}

func NewHost() *Receiver {
	return &Receiver{
		chDone: make(chan struct{}),
	}
}

func (r Receiver) Listen() (string, error) {
	lookupListener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 1111})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer lookupListener.Close()
	// listener.SetDeadline(time.Now().Add(2 * time.Minute))
	lookupConn, err := lookupListener.AcceptTCP()
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer lookupConn.Close()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 1112})
	if err != nil {
		log.Println(err)
		return "", err
	}
	// listener.SetDeadline(time.Now().Add(2 * time.Minute))
	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Println(err)
		return "", err
	}
	r.listener = listener
	r.connection = conn
	go r.receiveAsync()
	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	return ip, nil
}

func (r Receiver) Close() error {
	close(r.chDone)
	return errors.Join(r.connection.Close(), r.listener.Close())
}

func (r Receiver) receiveAsync() {
	for {
		select {
		case <-r.chDone:
			return

		case <-time.After(time.Second):
			if r.files == nil {
				err := gob.NewDecoder(r.connection).Decode(&r.files)
				if err != nil {
					runtime.EventsEmit(appCtx, "communication", "failed", "", err.Error())
					r.Close()
				}
			}
			for _, file := range r.files {
				runtime.EventsEmit(appCtx, "receiving", "preparing", file.Path, file.Size)

				// chBytesWrote := make(chan int64)
				// go func() {
				// 	for syncedBytes := range chBytesWrote {
				// 		totalSyncedBytes += syncedBytes
				// 		donePercent := int(syncedBytes * 100 / file.Size)
				// 		runtime.EventsEmit(appCtx, "receiving", "inProgress", file.Path, donePercent)
				// 	}
				// }()

				// --new
				//
				// err := tmpWrite(file, h.connection, file.Size)
				//
				// --end new
				// r := io.LimitReader(h.connection, file.Size)
				// err := writeFile(file, r, config.DefaultBufferSize, chBytesWrote)
				_, err := newWrite(file, r.connection, func(_ int64, percent int) {
					runtime.EventsEmit(appCtx, "receiving", "inProgress", file.Path, percent)
				})
				if err != nil {
					runtime.EventsEmit(appCtx, "receiving", "failed", file.Path, err)
					continue
				}
				runtime.EventsEmit(appCtx, "receiving", "completed", file.Path, "")
			}
			r.files = nil
		}
	}
}

func newWrite(file config.File, r io.Reader, onWrite func(transferred int64, percent int)) (int64, error) {
	filePath := filepath.Join(file.Parents...)
	filePath = filepath.Join(config.DefaultDirectory(), filePath, file.Name)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return 0, err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	tr := newTread(file.Size, onWrite)
	tee := io.TeeReader(r, tr)
	return io.CopyN(f, tee, file.Size)
}

func tmpWrite(file config.File, r io.Reader, size int64) error {
	filePath := filepath.Join(file.Parents...)
	filePath = filepath.Join(config.DefaultDirectory(), filePath, file.Name)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.CopyN(f, r, size)
	return err
}

func writeFile(file config.File, r io.Reader, bufferSize int, chBytesWrote chan<- int64) error {
	filePath := filepath.Join(file.Parents...)
	filePath = filepath.Join(config.DefaultDirectory(), filePath, file.Name)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return read(r, f, bufferSize, chBytesWrote)
}
