package gui

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/amirzayi/relay/config"
	"github.com/amirzayi/relay/pkg/fileutil"
	"github.com/amirzayi/relay/pkg/netutil"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Sender struct {
	conn    net.Conn
	chFiles chan fileutil.Files
	chFile  chan fileutil.File
	chDone  chan struct{}
}

func NewSender() *Sender {
	return &Sender{
		chFiles: make(chan fileutil.Files),
		chFile:  make(chan fileutil.File),
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
			address := net.JoinHostPort(ip.String(), strconv.Itoa(config.LookupPort))
			conn, err := net.DialTimeout("tcp", address, config.LookupTimeout)
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
	address := net.JoinHostPort(ip.String(), strconv.Itoa(config.DefaultPort))
	conn, err := net.DialTimeout("tcp", address, config.DefaultTimeout)
	if err != nil {
		return err
	}
	s.conn = conn
	go s.consumeFiles()
	go s.sendAsync()
	return nil
}

func (s Sender) addFiles(filePaths ...string) (fileutil.Files, error) {
	files, err := fileutil.GetFilesByPaths(filePaths...)
	if err != nil {
		return nil, err
	}
	s.chFiles <- files
	return files, nil
}

func (s Sender) SendFiles() (fileutil.Files, error) {
	filePaths, err := runtime.OpenMultipleFilesDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return nil, err
	}
	return s.addFiles(filePaths...)
}

func (s Sender) SendDirectory() (fileutil.Files, error) {
	selection, err := runtime.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return nil, err
	}
	if selection == "" {
		return fileutil.Files{}, nil
	}
	return s.addFiles(selection)
}

func (s Sender) consumeFiles() {
	for {
		select {
		case <-s.chDone:
			return
		case files := <-s.chFiles:
			if err := files.SendDetails(s.conn); err != nil {
				runtime.EventsEmit(appCtx, "communication", "failed", "", err.Error())
				s.Close()
				return
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
			_, err := file.ProgressiveRead(s.conn, func(_ int64, percent int) {
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
