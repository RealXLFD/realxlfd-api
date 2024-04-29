package db

import (
	"database/sql"
	"os"

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
VALUES (?, ?, ?, ?);`, data.Hash, data.Main, data.Scale, data.Date,
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
			data.Main,
			data.Date,
		),
	)
	return
Error:
	log.Error(
		str.T(
			"sql error: can not add image({}, {}, {}) into table Images",
			data.Hash,
			data.Main,
			data.Date,
		),
	)
}
func (db *Sqlite) AddImageData(data *ImageData) {
	stat, err := os.Stat(data.Path)
	if err != nil || stat.IsDir() {
		log.Error(str.T("can not get content size of image: {path}", data.Path))
		return
	}
	contentSize := stat.Size()
	d := db.driver
	var affected int64
	result, err := d.Exec(
		`INSERT OR IGNORE INTO ImageData (Path, Hash, Size, Quality, Format,ContentSize)
VALUES (?, ?, ?, ?, ?, ?);`, data.Path, data.Hash, data.Size, data.Quality,
		data.Format, contentSize,
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
