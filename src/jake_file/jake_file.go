package jake_file

import (
	"os"
	"syscall"
)

type Jake_File struct {
	fd   int    // file descriptor number
	name string // file name at Open time
}

func newFile(fd int, name string) *Jake_File {
	if fd < 0 {
		return nil
	}
	return &Jake_File{fd, name}
}

var (
	Stdin  = newFile(0, "/dev/stdin")
	Stdout = newFile(1, "/dev/stdout")
	Stderr = newFile(2, "/dev/stderr")
)

func Open(name string, mode int, perm uint32) (file *Jake_File, err error) {
	r, e := syscall.Open(name, mode, perm)
	if e != nil {
		err = e
	}
	return newFile(r, name), err
}

func (file *Jake_File) Close() error {
	if file == nil {
		return os.ErrInvalid
	}
	e := syscall.Close(file.fd)
	file.fd = -1 // so it can't be closed again
	if e != nil {
		return e
	}
	return nil
}

func (file *Jake_File) Read(b []byte) (ret int, err error) {
	if file == nil {
		return -1, os.ErrInvalid
	}
	r, e := syscall.Read(file.fd, b)
	if e != nil {
		err = e
	}
	return int(r), err
}

func (file *Jake_File) Write(b []byte) (ret int, err error) {
	if file == nil {
		return -1, os.ErrInvalid
	}
	r, e := syscall.Write(file.fd, b)
	if e != nil {
		err = e
	}
	return int(r), err
}

func (file *Jake_File) String() string {
	return file.name
}
