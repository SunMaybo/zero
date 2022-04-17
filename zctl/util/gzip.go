package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func TarGz(filepath, filename string) error {
	File, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer File.Close()
	gw := gzip.NewWriter(File)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	return walk(filepath, tw)
}
func UnTarGz(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}
func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func walk(path string, tw *tar.Writer) error {
	path = strings.Replace(path, "\\", "/", -1)
	info, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	index := strings.Index(path, "/")
	list := strings.Join(strings.Split(path, "/")[index:], "/")
	for _, v := range info {
		if v.IsDir() {
			head := tar.Header{Name: list + v.Name(), Typeflag: tar.TypeDir, ModTime: v.ModTime()}
			tw.WriteHeader(&head)
			walk(path+v.Name(), tw)
			continue
		}
		F, err := os.Open(path + v.Name())
		if err != nil {
			fmt.Println("打开文件%s失败.", err)
			continue
		}
		head := tar.Header{Name: list + v.Name(), Size: v.Size(), Mode: int64(v.Mode()), ModTime: v.ModTime()}
		tw.WriteHeader(&head)
		io.Copy(tw, F)
		F.Close()
	}
	return nil
}
