package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	imageTableName    = "image"
	tagTableName      = "tag"
	imageTagTableName = "tag"
)

func initDB() {
	db, err := sql.Open("sqlite", cfg.DbName)
	logErrorIfExistsAndExit(err)

	_, execErr := db.Exec(`
	drop table if exists tag;
	create table tag(name type UNIQUE);
	drop table if exists image;
	create table image(hash type UNIQUE, name);
	drop table if exists image_tag;
	create table image_tag(image_id, tag_id);
	`)
	logErrorIfExistsAndExit(execErr)

	db.Close()
}

func openDB() *sql.DB {
	db, err := sql.Open("sqlite", cfg.DbName)
	if err != nil {
		slog.Error(err.Error())
	}
	return db
}

func dbCreateTagEntry(db *sql.DB, tagName string) {
	_, err := db.Exec("insert into tag values(?);", tagName)
	logErrorIfExists(err)
}

func dbAddTag(db *sql.DB, image id, tag id) {
	_, err := db.Exec("insert into image_tag values(?, ?);", image, tag)
	logErrorIfExists(err)
}

func dbCreateFileEntry(db *sql.DB, hash byteHash, name string) {
	_, err := db.Exec("insert into image values(?, ?);", hash, name)
	logErrorIfExists(err)
}

func dbFetchTagsByNames(db *sql.DB, tagNames []string) []dbTag {
	addQuotation := func(x string) string { return fmt.Sprintf("'%s'", x) }
	var tags []dbTag
	statement := fmt.Sprintf("SELECT _rowid_, name FROM tag WHERE name IN (%s)", strings.Join(mapping(tagNames, addQuotation), ","))
	rows, err := db.Query(statement)
	logErrorIfExists(err)
	defer rows.Close()
	for rows.Next() {
		var (
			id   int64
			name string
		)
		err := rows.Scan(&id, &name)
		logErrorIfExists(err)
		tags = append(tags, *newDbTag(id, name))
	}
	err = rows.Err()
	logErrorIfExists(err)
	return tags
}

func dbFetchTagByName(db *sql.DB, tagName string) dbTag {
	row := db.QueryRow("SELECT _rowid_, name FROM tag WHERE name = (?)", tagName)
	var (
		id   int64
		name string
	)
	err := row.Scan(&id, &name)
	logErrorIfExists(err)
	return *newDbTag(id, name)
}

func dbFetchImagesByTags(db *sql.DB, tags []dbTag) []dbImage {
	getIds := func(x dbTag) string { return fmt.Sprint(x.id) }
	var images []dbImage
	statement := fmt.Sprintf("SELECT image._rowid_, hash, image.name FROM image INNER JOIN image_tag ON image._rowid_ = image_id INNER JOIN tag ON tag_id = tag._rowid_ WHERE tag_id IN (%s)", strings.Join(mapping(tags, getIds), ","))
	rows, err := db.Query(statement)
	logErrorIfExists(err)
	defer rows.Close()
	for rows.Next() {
		var (
			id   int64
			hash byteHash
			name string
		)
		err := rows.Scan(&id, &hash, &name)
		logErrorIfExists(err)
		images = append(images, *newDbImage(id, toHexHash(hash), name))
	}
	err = rows.Err()
	logErrorIfExists(err)
	return images
}

func dbFetchImages(db *sql.DB) []dbImage {
	var images []dbImage
	rows, err := db.Query("SELECT image._rowid_, hash, image.name FROM image")
	logErrorIfExists(err)
	defer rows.Close()
	for rows.Next() {
		var (
			id   int64
			hash byteHash
			name string
		)
		err := rows.Scan(&id, &hash, &name)
		logErrorIfExists(err)
		images = append(images, *newDbImage(id, toHexHash(hash), name))
	}
	err = rows.Err()
	logErrorIfExists(err)
	return images
}

func dbFetchImageById(db *sql.DB, queryId id) dbImage {
	row := db.QueryRow("SELECT _rowid_, hash, name FROM image WHERE _rowid_ = (?)", queryId)
	var (
		id   int64
		hash byteHash
		name string
	)
	err := row.Scan(&id, &hash, &name)
	logErrorIfExists(err)
	return *newDbImage(id, toHexHash(hash), name)
}

func dbFetchTagsForImage(db *sql.DB, imageId id) []dbTag {
	var tags []dbTag
	rows, err := db.Query("SELECT tag._rowid_, tag.name FROM image INNER JOIN image_tag ON image._rowid_ = image_id INNER JOIN tag ON tag_id = tag._rowid_ WHERE image._rowid_ = (?)", imageId)
	logErrorIfExists(err)
	defer rows.Close()
	for rows.Next() {
		var (
			id   int64
			name string
		)
		err := rows.Scan(&id, &name)
		logErrorIfExists(err)
		tags = append(tags, *newDbTag(id, name))
	}
	err = rows.Err()
	logErrorIfExists(err)
	return tags
}
