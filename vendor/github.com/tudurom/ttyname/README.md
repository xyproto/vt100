# ttyname

[![GoDoc](https://godoc.org/github.com/tudurom/ttyname?status.svg)](https://godoc.org/github.com/tudurom/ttyname)

Prints the file name of the terminal connected to standard input.

## Usage

```go
tty, err := ttyname.TTY()
if err != nil {
	panic(err)
}

fmt.Println(tty)
```
