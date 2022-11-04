package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"share/internal/file_system"
	"share/internal/service"
	"strings"
	"syscall"
)

var (
	//torrentMode = flag.Bool("torrent", false, "path")
	help   = flag.Bool("help", false, "show help")
	path   = flag.String("path", "", "path to directory or file. if empty - current directory will be used.")
	ignore = flag.String("ignore", "", "comma separated list of files or directories to ignore")
	only   = flag.String("only", "", "comma separated list of files or directories to serve")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	targetPath := *path
	if targetPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("can't get current directory: %s", err)
		}
		targetPath = pwd
	}

	fs := file_system.NewFileSystem(targetPath, getList(*ignore), getList(*only))
	server := service.NewFileServer(fs)
	go shutdown(server)
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("can't start file server: %s", err)
	}
}

func getList(str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
}

func shutdown(server *service.FileServer) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	if err := server.Stop(); err != nil {
		log.Fatalf("can't stop file server: %s", err)
	}
}
