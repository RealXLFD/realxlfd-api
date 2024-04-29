package db

import (
	"database/sql"
	"errors"
	"strconv"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

func (db *Sqlite) Count(req *RpicRequest) (count int, ok bool) {
	template := "SELECT Count(DISTINCT Images.Hash) AS count FROM Images JOIN Albums ON Images.Hash = Albums.Hash {};"
	var result *sql.Row
	scaleQuery := ParseScale(req.Scale)
	switch {
	case req.Album != "" && scaleQuery != "":
		query := str.T(template, "WHERE Albums.Album = ? AND Images.Scale = ?")
		result = db.driver.QueryRow(query, req.Album, req.Scale)
	case req.Album != "":
		query := str.T(template, "WHERE Albums.Album = ?")
		result = db.driver.QueryRow(query, req.Album)
	case scaleQuery != "":
		query := str.T(template, "WHERE Images.Scale = ?")
		result = db.driver.QueryRow(query, req.Scale)
	default:
		query := str.T(template, "")
		result = db.driver.QueryRow(query)
	}
	err := result.Scan(&count)
	if err != nil {
		log.Error(
			str.T(
				"sql error: can not get the image count in album({}) with scale({})",
				req.Album,
				req.Scale,
			),
		)
		return 0, false
	}
	log.Debug(
		str.T(
			"find {} image in album({}) with scale({})", count, req.Album,
			req.Scale,
		),
	)
	return count, true
}

// GetPath alert: ok=true时也包括Path未找到的情况
func (db *Sqlite) GetPath(i *ImageData) (path string, contentSize int64, ok bool) {
	query := "SELECT Path,ContentSize FROM ImageData WHERE Hash = ? AND Format = ? AND Size = ? AND Quality = ?;"
	result := db.driver.QueryRow(query, i.Hash, i.Format, i.Size, i.Quality)
	err := result.Scan(&path, &contentSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug(
				str.T(
					"no path found with image({}, {}, {}, {})", i.Hash, i.Format, i.Size,
					strconv.Itoa(i.Quality),
				),
			)
			return "", 0, true
		}
		log.Error("sql error: can not execute query({}) in table ImageData", query)
		return "", 0, false
	}
	return path, contentSize, true
}

func (db *Sqlite) GetAllPaths(hash string) (paths []string, ok bool) {
	query := "SELECT Path FROM ImageData WHERE Hash = ?;"
	result, err := db.driver.Query(query, hash)
	if err != nil {
		log.Error("sql error: can not execute query({}) in table ImageData", query)
		return nil, false
	}
	for result.Next() {
		var path string
		err = result.Scan()
		if err != nil {
			log.Error(str.T("sql error: {}", err.Error()))
			return nil, false
		}
		paths = append(paths, path)
	}
	log.Debug(str.T("find {} paths for image({})", len(paths), hash))
	return paths, true
}

func (db *Sqlite) Contains(hash, album string) bool {
	var row *sql.Row
	if album == "" {
		query := "SELECT Hash FROM Images WHERE Hash = ?;"
		row = db.driver.QueryRow(query, hash)
	} else {
		query := "SELECT Hash FROM Albums WHERE Hash = ? AND Album = ?;"
		row = db.driver.QueryRow(query, hash, album)
	}
	err := row.Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	return true
}

func (db *Sqlite) CountAlbums(hash string) (count int, ok bool) {
	query := "SELECT Count(Album) FROM Albums WHERE Hash = ?;"
	row := db.driver.QueryRow(query, hash)
	err := row.Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false
	}
	return count, true
}
