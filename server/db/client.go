package db

import (
	"database/sql"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

type Sqlite struct {
	driver *sql.DB
}

func (db *Sqlite) AddAlbum(hash, album string) {
	d := db.driver
	var affected int64
	result, err := d.Exec(
		`INSERT OR IGNORE INTO Albums (Hash,Album) VALUES (?, ?);`, hash, album,
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
			"{} rows affected: add image({}) into Album({})", affected, hash,
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

func (db *Sqlite) Rpic(req *RpicRequest) {
	template := "SELECT {} FROM {} WHERE {} ORDER BY {} LIMIT {};"
	switch {
	case req.Album != "" && req.Scale != "":
		query := str.T(
			template, "Images.Hash",
			"Images JOIN Albums ON Images.Hash = Albums.Hash",
			"Albums.Album = ? AND Images.Scale = ?", "RANDOM()", "1",
		)
		result := db.driver.QueryRow(query, req.Album, req.Scale)
		_ = result.Scan("Hash")
	}

}
