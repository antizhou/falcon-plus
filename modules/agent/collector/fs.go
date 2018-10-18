package collector

import (
	"os"
	"time"

	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

func SetMemMapFs() {
	fs = afero.NewMemMapFs()
}

func Create(name string) (afero.File, error) {
	return fs.Create(name)
}

func Mkdir(name string, perm os.FileMode) error {
	return fs.Mkdir(name, perm)
}

func MkdirAll(path string, perm os.FileMode) error {
	return fs.MkdirAll(path, perm)
}

func Open(name string) (afero.File, error) {
	return fs.Open(name)
}

func OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return fs.OpenFile(name, flag, perm)
}

func Remove(name string) error {
	return fs.Remove(name)
}

func RemoveAll(path string) error {
	return fs.RemoveAll(path)
}

func Rename(oldname, newname string) error {
	return fs.Rename(oldname, newname)
}

func Stat(name string) (os.FileInfo, error) {
	return fs.Stat(name)
}

func Chmod(name string, mode os.FileMode) error {
	return fs.Chmod(name, mode)
}

func Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.Chtimes(name, atime, mtime)
}
