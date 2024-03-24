package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func calculateHashes(files []filePath) []fileHash {
	var hashes []fileHash
	numJobs := len(files)
	jobs := make(chan string, numJobs)
	results := make(chan fileHash, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(jobs, results)
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		hashes = append(hashes, <-results)
	}

	return hashes
}

func worker(jobs <-chan filePath, results chan<- fileHash) {
	for j := range jobs {
		results <- createSHA256(j)
	}
}

func createSHA256(filePath filePath) fileHash {
	f, err := os.Open(filePath)
	logErrorIfExists(err)
	defer f.Close()

	h := sha256.New()
	_, copyErr := io.Copy(h, f)
	logErrorIfExists(copyErr)

	return newFileHash(filePath, h.Sum(nil))
}

func retrieveFileByHash(cfg config, hash []byte) []byte {
	return readFile(getOutPath(cfg, hash))
}

func getOutPath(cfg config, hash []byte) filePath {
	return filepath.Join(cfg.OutDir, strings.Join([]string{hex.EncodeToString(hash), ".jpg"}, ""))
}
