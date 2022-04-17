package file

import (
	"os"
	"runtime"
	"strings"
)

func PathExists(path string) (bool, error) {
	path = getPath(path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriterFile(filePath string, data []byte) error {
	filePath = getPath(filePath)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func ReadFile(filePath string) ([]byte, error) {
	filePath = getPath(filePath)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data := make([]byte, 1024)
	n, err := f.Read(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}
func MkdirAll(dir string) error {
	dir = getPath(dir)
	return os.MkdirAll(dir, 0755)
}
func ChmodExecuteFile(filePath string) error {
	filePath = getPath(filePath)
	return os.Chmod(filePath, 0755)
}
func RemoveAll(dir string) error {
	dir = getPath(dir)
	return os.RemoveAll(dir)
}
func GetFilePath(baseDir, filePath string) string {
	if runtime.GOOS == "windows" {
		filePath = strings.ReplaceAll(filePath, "/", "\\")
		return baseDir + filePath
	} else {
		return baseDir + filePath
	}
}
func getPath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "/", "\\")
	}
	return path
}
