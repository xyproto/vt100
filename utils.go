package vt100

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// logf, for quick "printf-style" debugging
func logf(head string, tail ...interface{}) {
	tmpdir := os.Getenv("TMPDIR")
	if tmpdir == "" {
		tmpdir = "/tmp"
	}
	logfilename := filepath.Join(tmpdir, "o.log")
	f, err := os.OpenFile(logfilename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		f, err = os.Create(logfilename)
		if err != nil {
			log.Fatalln(err)
		}
	}
	f.WriteString(fmt.Sprintf(head, tail...))
	f.Sync()
	f.Close()
}

// Silence the "logf is unused" message by staticcheck
var _ = logf
