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
		`INSERT OR IGNORE INTO Images (Hash, Scale, Date)
VALUES (?, ?, ?);`, data.Hash, data.Scale, data.Date,
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
			"{} rows affected: add image({}, {}, {})", affected, data.Hash,
			data.Scale,
			data.Date,
		),
	)
	return
Error:
	log.Error(
		str.T(
			"sql error: can not add image({}, {}, {}) into table Images",
			data.Hash,
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

func (db *Sqlite) Rpic(req *RpicRequest) (hash string, ok bool) {
	template := "SELECT {} FROM {} {} ORDER BY {} LIMIT 1 {};"
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
			return "", false
		}
		if req.Rid < 0 {
			req.Rid = -req.Rid
		}
		dArgOffset = str.Join("OFFSET ", strconv.Itoa(req.Rid%count))
	}
	switch {
	case req.Album != "" && req.Scale != "":
		query := str.T(
			template, "Images.Hash",
			"Images JOIN Albums ON Images.Hash = Albums.Hash",
			"WHERE Albums.Album = ? AND Images.Scale = ?", dArgOrder, dArgOffset,
		)
		result = db.driver.QueryRow(query, req.Album, req.Scale)
	case req.Album != "":
		query := str.T(template, "Hash", "Albums", "WHERE Album = ?", dArgOrder, dArgOffset)
		result = db.driver.QueryRow(query, req.Album)
	case req.Scale != "":
		query := str.T(template, "Hash", "Images", "WHERE Scale = ?", dArgOrder, dArgOffset)
		result = db.driver.QueryRow(query, req.Scale)
	default:
		query := str.T(template, "Hash", "Images", "", dArgOrder, dArgOffset)
		result = db.driver.QueryRow(query)
	}
	err := result.Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			goto NotFound
		}
		goto Error
	}
	log.Debug(
		str.T(
			"get rpic({}) ", hash,
		),
	)
	return
NotFound:
	log.Debug(str.T("no image was found in album({}) with scale({})", req.Album, req.Scale))
Error:
	log.Error(
		str.T(
			"sql error: can not get random image in album({}) with scale({})",
			req.Album, req.Scale,
		),
	)
	return "", false
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
		return count, false
	}
	log.Info(str.T("find {} image in album({}) with scale({})"), req.Album, req.Scale)
	return count, true
}

func (db *Sqlite) GetPath(hash string,format ) {
	format size
}

func ()