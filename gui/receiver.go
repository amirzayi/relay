package gui

import (
	"encoding/gob"
	"errors"
	"net"
	"strings"

	"github.com/amirzayi/relay/config"
	"github.com/amirzayi/relay/pkg/fileutil"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Receiver struct {
	listener   *net.TCPListener
	connection *net.TCPConn
	files      fileutil.Files
	chDone     chan struct{}
}

func NewHost() *Receiver {
	return &Receiver{
		chDone: make(chan struct{}),
	}
}

func (r Receiver) Listen() (string, error) {
	lookupListener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.LookupPort})
	if err != nil {
		return "", err
	}
	defer lookupListener.Close()

	// listener.SetDeadline(time.Now().Add(2 * time.Minute))
	lookupConn, err := lookupListener.AcceptTCP()
	if err != nil {
		return "", err
	}
	defer lookupConn.Close()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: config.DefaultPort})
	if err != nil {
		return "", err
	}
	// listener.SetDeadline(time.Now().Add(2 * time.Minute))

	conn, err := listener.AcceptTCP()
	if err != nil {
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

		default:
			if r.files == nil {
				err := gob.NewDecoder(r.connection).Decode(&r.files)
				if err != nil {
					runtime.EventsEmit(appCtx, "communication", "failed", "", err.Error())
					r.Close()
					return
				}
			}
			for _, file := range r.files {
				runtime.EventsEmit(appCtx, "receiving", "preparing", file.Path, file.Size)
				_, err := file.ProgressiveWrite(r.connection, func(_ int64, percent int) {
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
