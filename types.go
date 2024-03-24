package main

type dirPath = string
type filePath = string
type message = string
type html = string
type base64String = string
type byteHash = []byte
type hexHash = string
type id = int64

type fileHash struct {
	path string
	hash byteHash
}

func newFileHash(path filePath, hash []byte) fileHash {
	return fileHash{path, hash}
}

type CopyInOut struct {
	in  string
	out string
}

func newCopyInOut(in filePath, out filePath) CopyInOut {
	return CopyInOut{in, out}
}

type messageRequest struct {
	msg  message
	body []byte
}

func newMessageRequest(msg message, body []byte) messageRequest {
	return messageRequest{msg, body}
}

type config struct {
	ConsumeDir string
	OutDir     string
	DbName     string
	Url        string
	Port       string
}

type image struct {
	Id      string
	Content base64String
}

func newImage(id string, content base64String) *image {
	return &image{id, content}
}

type dbTag struct {
	id   id
	name string
}

func newDbTag(id id, name string) *dbTag {
	return &dbTag{id, name}
}

type dbImage struct {
	id   id
	hash hexHash
	name string
}

func newDbImage(id id, hash hexHash, name string) *dbImage {
	return &dbImage{id, hash, name}
}

type uiImage struct {
	Id      id
	Content base64String
}

func newUiImage(id id, content base64String) *uiImage {
	return &uiImage{id, content}
}

type uiDetailImage struct {
	Id      id
	Content base64String
	Tags    []string
}

func newUiDetailImage(id id, content base64String, tags []string) *uiDetailImage {
	return &uiDetailImage{id, content, tags}
}
