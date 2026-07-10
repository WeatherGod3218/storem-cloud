package verify

import (
	"sync"
)

var filesToVerify []string = make([]string, 0)
var mutex sync.Mutex = sync.Mutex{}

func AddFileToVerifyList(file string) {
	mutex.Lock()
	defer mutex.Unlock()

	filesToVerify = append(filesToVerify, file)
}

func GetVerifyList() []string {
	return filesToVerify
}
