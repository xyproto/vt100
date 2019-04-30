// +build !windows

package ttyname

import "os"

const fd0 = "/proc/self/fd/0"

// TTY prints the file name of the terminal connected to standard input
func TTY() (string, error) {
	dest, err := os.Readlink(fd0)
	if err != nil {
		return "", err
	}
	return dest, nil
}
