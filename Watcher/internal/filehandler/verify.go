package filehandler

import (
	"sync"
)

var filesToVerify []string = make([]string, 0)
var mutex sync.Mutex = sync.Mutex{}

// const MAX_CONCURRENT_SCANS = 4

// var sem chan struct{} = make(chan struct{}, MAX_CONCURRENT_SCANS)

func AddFileToVerifyList(file string) {
	mutex.Lock()
	defer mutex.Unlock()

	filesToVerify = append(filesToVerify, file)
}

func GetHashedFilesList() []string {
	return filesToVerify
}
