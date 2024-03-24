package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cbroglie/mustache"
)

func index() html {
	styleBase64 := encodeStringBase64(style)
	htmxBase64 := encodeStringBase64(htmx)

	db := openDB()
	dbImages := dbFetchImages(db)

	uiImages := mapping(dbImages, func(x dbImage) uiImage {
		return *newUiImage(x.id, getContent(x))
	})

	return renderTemplate(indexHtml, map[string]interface{}{
		"heading": "Hello Heading",
		"style":   styleBase64,
		"htmx":    htmxBase64,
		"imgs":    uiImages})
}

func show() html {
	return "<p>show message send</p>"
}

func renderTemplate(template string, args map[string]interface{}) html {
	html, err := mustache.Render(template, args)
	logErrorIfExists(err)
	return html
}

func getImages(path dirPath) []image {
	var images []image
	files := readDir(path)
	for _, f := range files {
		content := readFile(f)
		id := strings.Split(filepath.Base(f), ".")[0]
		images = append(images, *newImage(id, encodeBytesBase64(content)))
	}

	return images
}

func showImage(imageId id) html {
	db := openDB()
	return renderUiDetailImage(db, imageId)
}

func addTag(imageId id, tagName string) html {
	db := openDB()
	dbCreateTagEntry(db, tagName)
	dbTag := dbFetchTagByName(db, tagName)
	dbAddTag(db, imageId, dbTag.id)
	dbFetchImageById(db, imageId)
	return renderUiDetailImage(db, imageId)
}

func renderUiDetailImage(db *sql.DB, imageId id) html {
	dbImage := dbFetchImageById(db, imageId)
	dbTags := dbFetchTagsForImage(db, dbImage.id)
	uiDetailImage := newUiDetailImage(dbImage.id, getContent(dbImage), mapping(dbTags, func(x dbTag) string { return x.name }))
	fmt.Println(uiDetailImage.Tags)
	return renderTemplate(modalHtml, map[string]interface{}{
		"img": uiDetailImage})
}

func searchImages(tagName string) html {
	db := openDB()
	tag := dbFetchTagByName(db, tagName)
	dbImages := dbFetchImagesByTags(db, []dbTag{tag})
	uiImages := mapping(dbImages, func(x dbImage) uiImage {
		return *newUiImage(x.id, getContent(x))
	})
	return renderTemplate(searchHtml, map[string]interface{}{
		"imgs": uiImages})
}

func getContent(image dbImage) base64String {
	return encodeBytesBase64(readFile(buildPath(strings.Join([]string{image.hash, filepath.Ext(image.name)}, ""))))
}
