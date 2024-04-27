package db

import (
	"database/sql"
	"errors"
	"strconv"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

type Sqlite struct {
	driver *sql.DB
}

func (db *Sqlite) AddAlbum(hash, album string) {
	d := db.driver
	var affected int64
	result, err := d.Exec(
		`INSERT OR IGNORE INTO Albums (Hash,Album) VALUES (?, ?);`,
		hash,
		album,
	)
	if err != nil {
		goto Error
	}
	affected, err = result.RowsAffected()
	if err != nil {
		goto Error
	}
	log.Debug(
		str.T(
			"{} rows affected: add image({}) into Album({})",
			affected,
			hash,
			album,
		),
	)
	return
Error:
	log.Error(
		str.T(
			"sql error: can not add image({}) into Album({})", hash, album,
		),
	)
}
func (db *Sqlite) AddImage(data *Image) {
	d := db.driver
	var affected int64
	result, err := d.Exec(
		`INSERT OR IGNORE INTO Images (Hash, Main, Scale, Date)
VALUES (?, ?, ?, ?);`, data.Hash, data.Scale, data.Date,
	)
	if err != nil {
		goto Error
	}
	affected, err = result.RowsAffected()
	if err != nil {
		goto Error
	}
	log.Debug(
		str.T(
			"{} rows affected: add image({}, {}, {}, {})", affected, data.Hash,
			data.Main,
			data.Scale,
			data.Date,
		),
	)
	return
Error:
	log.Error(
		str.T(
			"sql error: can not add image({}, {}, {}, {}) into table Images",
			data.Hash,
			data.Main,
			data.Scale,
			data.Date,
		),
	)
}
func (db *Sqlite) AddImageData(data *ImageData) {
	d := db.driver
	var affected int64
	result, err := d.Exec(
		`INSERT OR IGNORE INTO ImageData (Path, Hash, Size, Width, Height, Format)
VALUES (?, ?, ?, ?, ?, ?);`, data.Path, data.Hash, data.Width, data.Height,
		data.Format,
	)
	if err != nil {
		goto Error
	}
	affected, err = result.RowsAffected()
	if err != nil {
		goto Error
	}
	log.Debug(
		str.T(
			"{} rows affected: add image_data({})", affected, data.Path,
		),
	)
	return
Error:
	log.Error(
		str.T(
			"sql error: can not add image_data({}) into table ImageData",
			data.Path,
		),
	)
}

func (db *Sqlite) Rpic(req *RpicRequest) (hash string, main string, ok bool) {
	template := "SELECT Images.Hash, Images.Main FROM Images JOIN Albums ON Images.Hash = Albums.Hash {} ORDER BY {} LIMIT 1 {};"
	var dArgOrder, dArgOffset string
	var result *sql.Row
	if req.HasRid {
		dArgOrder = "Images.Date"
		dArgOffset = ""
	} else {
		dArgOrder = "RANDOM()"
		var count int
		count, ok = db.Count(req)
		if !ok {
			return "", "", false
		}
		if req.Rid < 0 {
			req.Rid = -req.Rid
		}
		dArgOffset = str.Join("OFFSET ", strconv.Itoa(req.Rid%count))
	}
	switch {
	case req.Album != "" && req.Scale != "":
		query := str.T(
			template,
			"WHERE Albums.Album = ? AND Images.Scale = ?", dArgOrder,
			dArgOffset,
		)
		result = db.driver.QueryRow(query, req.Album, req.Scale)
	case req.Album != "":
		query := str.T(
			template, "WHERE Albums.Album = ?", dArgOrder,
			dArgOffset,
		)
		result = db.driver.QueryRow(query, req.Album)
	case req.Scale != "":
		query := str.T(
			template, "WHERE Images.Scale = ?", dArgOrder,
			dArgOffset,
		)
		result = db.driver.QueryRow(query, req.Scale)
	default:
		query := str.T(template, "", dArgOrder, dArgOffset)
		result = db.driver.QueryRow(query)
	}
	err := result.Scan(&hash, &main)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			goto NotFound
		}
		goto Error
	}
	log.Debug(
		str.T(
			"get rpic({}, {}) ", hash, main,
		),
	)
	return hash, main, true
NotFound:
	log.Debug(
		str.T(
			"no image was found in album({}) with scale({})", req.Album,
			req.Scale,
		),
	)
Error:
	log.Error(
		str.T(
			"sql error: can not get random image in album({}) with scale({})",
			req.Album, req.Scale,
		),
	)
	return "", "", false
}

func (db *Sqlite) Count(req *RpicRequest) (count int, ok bool) {
	template := "SELECT Count(DISTINCT Images.Hash) AS count FROM Images JOIN Albums ON Images.Hash = Albums.Hash {};"
	var result *sql.Row
	switch {
	case req.Album != "" && req.Scale != "":
		query := str.T(template, "WHERE Albums.Album = ? AND Images.Scale = ?")
		result = db.driver.QueryRow(query, req.Album, req.Scale)
	case req.Album != "":
		query := str.T(template, "WHERE Albums.Album = ?")
		result = db.driver.QueryRow(query, req.Album)
	case req.Scale != "":
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
		str.T("find {} image in album({}) with scale({})"), req.Album,
		req.Scale,
	)
	return count, true
}

func (db *Sqlite) GetPath(hash, format, size string) (path string, ok bool) {
	query := "SELECT Path FROM ImageData WHERE Hash = ? AND format = ? AND size = ?;"
	result := db.driver.QueryRow(query, hash, format, size)
	err := result.Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("find {} path with image({}, {}, {},)", hash, format, size)
			return "", true
		}
		log.Error("sql error: can not execute query({}) in table ImageData", query)
		return "", false
	}
	return path, true
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
