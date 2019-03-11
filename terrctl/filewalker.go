package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedisct1/dlog"
)

// WalkedFile - A file found during traverseal
type WalkedFile struct {
	Path         string
	RelativePath string
}

type fileWalker struct {
	files             []WalkedFile
	rootWithSeparator string
}

func (walker *fileWalker) visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fileName := filepath.Base(path)
	if strings.EqualFold(fileName, "node_modules") || strings.HasPrefix(fileName, ".") {
		dlog.Infof("Skipping [%v]", path)
		return filepath.SkipDir
	}
	if !f.Mode().IsRegular() {
		return nil
	}
	walkedFile := WalkedFile{
		Path:         path,
		RelativePath: strings.TrimPrefix(path, walker.rootWithSeparator),
	}
	dlog.Debugf("File [%v] added to the archive", walkedFile.RelativePath)
	walker.files = append(walker.files, walkedFile)
	return nil
}

// FileWalk - Traverse a directory and returns the list of files it contains
func FileWalk(root string) ([]WalkedFile, error) {
	root = filepath.Clean(root)
	if strings.HasSuffix(root, string(filepath.Separator)) {
		return nil, errors.New("The root directory cannot be uploaded")
	}
	rootWithSeparator := root + string(filepath.Separator)
	walker := fileWalker{rootWithSeparator: rootWithSeparator}
	err := filepath.Walk(root, walker.visit)
	if err != nil {
		return nil, err
	}
	return walker.files, nil
}
