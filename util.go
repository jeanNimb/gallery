package main

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"log/slog"
	"path/filepath"
)

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func mapping[T any, O any](ss []T, f func(T) O) (ret []O) {
	for _, s := range ss {
		ret = append(ret, f(s))
	}
	return
}

func encodeStringBase64(str string) base64String {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func encodeBytesBase64(b []byte) base64String {
	return base64.StdEncoding.EncodeToString([]byte(b))
}

func toHexHash(hash byteHash) hexHash {
	return hex.EncodeToString(hash)
}

func toByteHash(hash hexHash) byteHash {
	h, err := hex.DecodeString(hash)
	logErrorIfExists(err)
	return h
}

func logErrorIfExists(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

func logErrorIfExistsAndExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func buildPath(fileName string) filePath {
	return filepath.Join(cfg.OutDir, fileName)
}
