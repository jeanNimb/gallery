package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	_ "embed"
)

const configFile = "config.yaml"

//go:embed index.html
var indexHtml string

//go:embed modal.html
var modalHtml string

//go:embed search.html
var searchHtml string

//go:embed style.css
var style string

//go:embed htmx.min.js
var htmx string

var cfg config

type mode int

const (
	test mode = iota
	dev
)

func (m mode) String() string {
	return [...]string{"test", "dev"}[m]
}

func main() {
	m := dev
	cfg = loadConfig(configFile)

	switch m {
	case test:
		testDB()
	default:
		run()
	}
}

func run() {
	// TODO wenn prod einkommentieren
	// create new DB if not existent (delete old data in DB)
	// if _, err := os.Stat(cfg.DbName); errors.Is(err, os.ErrNotExist) {
	// 	initDB()
	// }
	initDB()

	// create dir needed structure
	createDirIfNotExists(cfg.ConsumeDir)
	createDirIfNotExists(cfg.OutDir)

	consumeFiles()

	// TODO only for testing
	db := openDB()
	dbCreateTagEntry(db, "test")
	images := dbFetchImages(db)
	for _, image := range images {
		dbAddTag(db, image.id, 1)
	}

	http.HandleFunc("/", handleRequest)
	fmt.Printf("server listen on %s:%s\n", cfg.Url, cfg.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Url, cfg.Port), nil)
}

func consumeFiles() {
	// open DB
	db := openDB()

	// read files in consume
	files := readDir(cfg.ConsumeDir)

	// calculate sha256 for each file
	hashes := calculateHashes(files)

	// copy files to out dir
	copyToOutDir(hashes, cfg.OutDir)

	// insert db entries for files
	for _, hash := range hashes {
		fileName := filepath.Base(hash.path)
		dbCreateFileEntry(db, hash.hash, fileName)
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	html := handleMessageRequest(newMessageRequest(getMessage(req), getBody(req)))
	fmt.Fprintf(w, "%s\n", html)
}

func handleMessageRequest(req messageRequest) html {
	switch req.msg {
	case "show":
		return show()
	case "show_img":
		bodyParams := getBodyParams(req)
		fmt.Printf("image id: %s\n", getBodyParams(req))
		imageId, err := strconv.ParseInt(bodyParams["Id"], 10, 64)
		logErrorIfExists(err)
		return showImage(imageId)
	case "close_modal":
		return ""
	case "search":
		bodyParams := getBodyParams(req)
		fmt.Printf("search body: %s\n", bodyParams)
		return searchImages(bodyParams["q"])
	case "add_tag":
		bodyParams := getBodyParams(req)
		fmt.Printf("add tag body: %s\n", bodyParams)
		imageId, err := strconv.ParseInt(bodyParams["Id"], 10, 64)
		logErrorIfExists(err)
		return addTag(imageId, bodyParams["tag"])
	default:
		return index()
	}
}

func getMessage(req *http.Request) message {
	return req.URL.Query().Get("msg")
}

func getBody(req *http.Request) []byte {
	b, err := io.ReadAll(req.Body)
	logErrorIfExists(err)

	return b
}

func getBodyParams(req messageRequest) map[string]string {
	bodyParams := make(map[string]string)
	for _, param := range strings.Split(string(req.body), "&") {
		keyValue := strings.Split(param, "=")
		bodyParams[keyValue[0]] = keyValue[1]
	}
	return bodyParams
}

func loadConfig(configFile filePath) config {
	var config config
	err := yaml.Unmarshal(readFile(configFile), &config)
	logErrorIfExistsAndExit(err)
	return config
}

func testDB() {
	initDB()
	db := openDB()
	dbCreateTagEntry(db, "image")
	dbCreateTagEntry(db, "test")
	dbCreateFileEntry(db, toByteHash("9afcef1fbf89b2ce37df878a67ce469a03f11d0d88d4a7e3ad1c15aa5c5ad008"), "dragon_1.jpg")
	dbCreateFileEntry(db, toByteHash("3b4c8f1fbf89b2ce37df878a67ce469a03f11d0d88d4a7e3ad1c15aa5c5ad008"), "dragon_2.jpg")
	dbAddTag(db, 1, 1)
	dbAddTag(db, 2, 1)
	tags := dbFetchTagsByNames(db, []string{"test"})
	fmt.Println(tags)
	images := dbFetchImagesByTags(db, []dbTag{*newDbTag(1, "image")})
	fmt.Println(images)
	paths := mapping(images, func(x dbImage) string { return buildPath(strings.Join([]string{x.hash, filepath.Ext(x.name)}, "")) })
	fmt.Println(paths)
}
