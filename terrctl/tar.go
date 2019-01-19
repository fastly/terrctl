package main

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
)

// CreateTarFile - Create a new Tar archive from a list of files
func CreateTarFile(files []WalkedFile) ([]byte, error) {
	var tarData bytes.Buffer
	tw := tar.NewWriter(&tarData)
	for _, file := range files {
		fd, err := os.Open(file.Path)
		if err != nil {
			return nil, err
		}
		stat, err := fd.Stat()
		if err != nil {
			return nil, err
		}
		hdr := &tar.Header{
			Name:    file.RelativePath,
			Mode:    0600,
			Size:    stat.Size(),
			ModTime: stat.ModTime(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		_, err = io.Copy(tw, fd)
		if err != nil {
			return nil, err
		}
	}
	tw.Close()
	return tarData.Bytes(), nil
}
