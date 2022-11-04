package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"share/internal/file_system"
	"share/internal/utils"
)

type FileServer struct {
	fileSystem *file_system.FileSystem
	httpServer *http.Server
}

func NewFileServer(fileSystem *file_system.FileSystem) *FileServer {
	return &FileServer{
		httpServer: &http.Server{},
		fileSystem: fileSystem,
	}
}

func (s *FileServer) Start() error {
	fs := http.FileServer(s.fileSystem)
	http.Handle("/", fs)
	listener, err := utils.NewListener()
	if err != nil {
		log.Fatalf("can't create listener: %s", err)
	}
	port, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		log.Fatalf("can't get port from listener")
	}
	localAddrList, err := utils.GetInterface()
	if err != nil {
		log.Fatalf("can't get local addr list: %s", err)
	}
	for _, localAddr := range localAddrList {
		url := fmt.Sprintf("http://%s:%d", localAddr, port.Port)
		log.Print("network address: " + url)
		if err = utils.PrintQRCode(url); err != nil {
			log.Fatalf("can't print QR code: %s", err)
		}
	}
	return s.httpServer.Serve(listener)
}

func (s *FileServer) Stop() error {
	return s.httpServer.Shutdown(context.Background())
}
