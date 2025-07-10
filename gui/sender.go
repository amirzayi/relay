package gui

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/amirzayi/relay/config"
	"github.com/amirzayi/relay/pkg/fileutil"
	"github.com/amirzayi/relay/pkg/netutil"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Sender struct {
	conn    net.Conn
	chFiles chan config.Files
	chFile  chan config.File
	chDone  chan struct{}
}

func NewSender() *Sender {
	return &Sender{
		chFiles: make(chan config.Files),
		chFile:  make(chan config.File),
		chDone:  make(chan struct{}),
	}
}

func (s Sender) Close() error {
	close(s.chDone)
	close(s.chFile)
	close(s.chFiles)
	return s.conn.Close()
}

func (s Sender) Lookup() ([]net.IP, error) {
	localIPs, err := netutil.GetLocalIPs()
	if err != nil {
		return nil, err
	}

	chIPs := make(chan net.IP)
	var wg sync.WaitGroup

	for _, localIP := range localIPs {
		wg.Add(1)
		go s.availableHosts(localIP, chIPs, &wg)
	}

	go func() {
		wg.Wait()
		close(chIPs)
	}()

	ips := []net.IP{}
	for ip := range chIPs {
		ips = append(ips, ip)
	}

	return ips, nil
}

func (s Sender) availableHosts(localIP net.IP, chIPs chan<- net.IP, wg *sync.WaitGroup) error {
	defer wg.Done()

	cidr := fmt.Sprintf("%s/24", localIP)
	ips, err := netutil.GetCIDRHosts(cidr)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()
			address := net.JoinHostPort(ip.String(), "1111")
			conn, err := net.DialTimeout("tcp", address, config.DefaultGUITimeout)
			if err != nil {
				return
			}
			chIPs <- ip
			conn.Close()
		}(ip)
	}
	return nil
}

func (s Sender) Connect(ip net.IP) error {
	address := net.JoinHostPort(ip.String(), config.DefaultGUIPort)
	conn, err := net.DialTimeout("tcp", address, config.DefaultGUITimeout)
	if err != nil {
		return err
	}
	s.conn = conn
	go s.consumeFiles()
	go s.sendAsync()
	return nil
}

func (s Sender) addFiles(filePaths ...string) (config.Files, error) {
	files, err := fileutil.GetFilesByPaths(filePaths...)
	if err != nil {
		return nil, err
	}
	s.chFiles <- files
	return files, nil
}

func (s Sender) SendFiles() (config.Files, error) {
	filePaths, err := runtime.OpenMultipleFilesDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return nil, err
	}
	return s.addFiles(filePaths...)
}

func (s Sender) SendDirectory() (config.Files, error) {
	selection, err := runtime.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return nil, err
	}
	if selection == "" {
		return config.Files{}, nil
	}
	return s.addFiles(selection)
}

func (s Sender) consumeFiles() {
	for {
		select {
		case <-s.chDone:
			return
		case files := <-s.chFiles:
			err := gob.NewEncoder(s.conn).Encode(files)
			if err != nil {
				runtime.EventsEmit(appCtx, "communication", "failed", "", err.Error())
				s.Close()
			}
			for _, file := range files {
				s.chFile <- file
			}
		}
	}
}

func (s Sender) sendAsync() {
	for {
		select {
		case <-s.chDone:
			return

		case file := <-s.chFile:
			_, err := newReader(file, s.conn, func(_ int64, percent int) {
				runtime.EventsEmit(appCtx, "sending", "inProgress", file.Path, percent)
			})
			if err != nil {
				runtime.EventsEmit(appCtx, "sending", "failed", file.Path, err)
				continue
			}
			runtime.EventsEmit(appCtx, "sending", "completed", file.Path, "")
		}
	}
}
func newReader(f config.File, w io.Writer, onWrite func(transferred int64, percent int)) (int64, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	t := newTread(f.Size, onWrite)
	r := io.TeeReader(file, t)
	return io.Copy(w, r)
}

func tmpRead(f config.File, w io.Writer, size int64) error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.CopyN(w, file, size)
	return err
}

type tread struct {
	totalSize        int64
	transferredBytes int64
	onWrite          func(transferred int64, percent int)
}

func newTread(totalSize int64, onWrite func(transferred int64, percent int)) io.Writer {
	return &tread{totalSize: totalSize, onWrite: onWrite}
}
func (t *tread) Write(p []byte) (n int, err error) {
	n = len(p)
	t.transferredBytes += int64(n)
	if t.onWrite != nil {
		size := t.totalSize
		if size == 0 {
			size = 1
		}
		t.onWrite(t.transferredBytes, int((t.transferredBytes*100)/size))
	}
	return
}

func readFile(f config.File, w io.Writer, bufferSize int, chBytesRead chan<- int64) error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	return read(file, w, bufferSize, chBytesRead)
}

func read(r io.Reader, w io.Writer, bufferSize int, chBytesRead chan<- int64) error {
	defer close(chBytesRead)

	rw := io.TeeReader(r, w)
	buffer := make([]byte, bufferSize)
	syncedBytes := int64(0)
	for {
		n, err := rw.Read(buffer)
		syncedBytes += int64(n)
		chBytesRead <- syncedBytes
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}
