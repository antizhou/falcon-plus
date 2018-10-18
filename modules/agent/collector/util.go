package collector

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func DirExists(path string) (bool, error) {
	return afero.DirExists(fs, path)
}

func Exists(path string) (bool, error) {
	return afero.Exists(fs, path)
}

func FileContainsBytes(filename string, subslice []byte) (bool, error) {
	return afero.FileContainsBytes(fs, filename, subslice)
}

func GetTempDir(subPath string) string {
	return afero.GetTempDir(fs, subPath)
}

func IsDir(path string) (bool, error) {
	return afero.IsDir(fs, path)
}

func IsEmpty(path string) (bool, error) {
	return afero.IsEmpty(fs, path)
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(fs, dirname)
}

func ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(fs, filename)
}

func SafeWriteReader(path string, r io.Reader) (err error) {
	return afero.SafeWriteReader(fs, path, r)
}

func TempDir(dir, prefix string) (name string, err error) {
	return afero.TempDir(fs, dir, prefix)
}

func TempFile(dir, prefix string) (f afero.File, err error) {
	return afero.TempFile(fs, dir, prefix)
}

func Walk(root string, walkFn filepath.WalkFunc) error {
	return afero.Walk(fs, root, walkFn)
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(fs, filename, data, perm)
}

func WriteReader(path string, r io.Reader) (err error) {
	return afero.WriteReader(fs, path, r)
}
