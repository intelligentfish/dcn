package io

import "os"

// PathExists check if path exits
func PathExists(path string) bool {
	fi, err := os.Stat(path)
	if nil != err && !os.IsNotExist(err) {
		panic(err)
	}
	return nil != fi
}
