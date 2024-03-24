package main

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

func copyToOutDir(hashes []fileHash, outDir dirPath) {
	getCopyInOut := func(fileHash fileHash) CopyInOut {
		src := fileHash.path
		dest := filepath.Join(outDir, strings.Join([]string{hex.EncodeToString(fileHash.hash), filepath.Ext(fileHash.path)}, ""))
		return newCopyInOut(src, dest)
	}

	copyInOutList := mapping(hashes, getCopyInOut)

	copyFiles(copyInOutList)
}

func readDir(path dirPath) []filePath {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		logErrorIfExists(err)
		files = append(files, path)
		return nil
	})
	logErrorIfExists(err)

	files = filter(files, func(s string) bool { return s != path })

	return files
}

func copyFiles(copyInOut []CopyInOut) {
	for _, c := range copyInOut {
		copy(c.in, c.out)
	}
}

func copy(src filePath, dst filePath) {
	data := readFile(src)
	writeFile(data, dst)
}

func readFile(src filePath) []byte {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	logErrorIfExists(err)
	return data
}

func writeFile(data []byte, dst filePath) {
	// Write data to dst
	err := os.WriteFile(dst, data, 0644)
	logErrorIfExists(err)
}

func createDirIfNotExists(path filePath) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		logErrorIfExists(err)
	}
}

func deleteDir(path filePath) {
	err := os.RemoveAll(path)
	logErrorIfExists(err)
}
