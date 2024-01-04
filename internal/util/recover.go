package util

import (
	"log"
)

func Recover(f func()) (panicked bool) {
	defer func() { // "recover" only works inside "defer".
		if r := recover(); r != nil {
			log.Println(r)
			panicked = true
		}
	}()

	f()

	return
}
