package file_system

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type FileSystem struct {
	fs         http.FileSystem
	path       string
	ignoreList []string
	onlyList   []string
}

func NewFileSystem(filePath string, ignoreList, onlyList []string) *FileSystem {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("path `%s` does not exist", filePath)
	}
	if len(ignoreList) > 0 {
		log.Printf("ignoring files: %v", ignoreList)
	}
	if len(onlyList) > 0 {
		log.Printf("serving only files: %v", onlyList)
	}
	log.Printf("starting file server for path: %s", filePath)
	return &FileSystem{
		fs:         http.Dir(filePath),
		onlyList:   onlyList,
		ignoreList: ignoreList,
		path:       filePath,
	}
}

func (nfs FileSystem) Open(path string) (http.File, error) {
	if !nfs.filterFile(path) {
		return nil, os.ErrNotExist
	}
	return nfs.fs.Open(path)
}

func (nfs FileSystem) filterFile(path string) bool {
	if len(nfs.onlyList) > 0 {
		return isInList(path, nfs.onlyList)
	}
	return !isInList(path, nfs.ignoreList)
}

func isInList(path string, list []string) bool {
	for _, item := range list {
		if strings.Contains(path, item) {
			return true
		}
		if strings.HasPrefix(item, "*") {
			// file mask here
			if strings.HasSuffix(path, item[1:]) {
				return true
			}
		}
	}
	return false
}
