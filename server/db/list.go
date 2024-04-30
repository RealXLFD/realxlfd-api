package db

import (
	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

func (db *Sqlite) ListAlbums() (albums []string) {
	query := "SELECT DISTINCT Album FROM Albums;"
	results, err := db.driver.Query(query)
	if err != nil {
		log.Error("sql failed: can not get albums from database")
		return nil
	}
	for results.Next() {
		var album string
		err = results.Scan(&album)
		if err != nil {
			log.Error("sql failed: can not get albums from database")
			return nil
		}
		albums = append(albums, album)
	}
	log.Debug(str.T("find {} albums in database", len(albums)))
	return
}

func (db *Sqlite) CountPics(album string) (count int, ok bool) {
	query := "SELECT COUNT(DISTINCT Hash) FROM Albums WHERE Album = ?;"
	row := db.driver.QueryRow(query, album)
	err := row.Scan(&count)
	if err != nil {
		log.Error(str.T("sql failed: can not count image id from album({})", album))
		return 0, false
	}
	return count, true
}

func (db *Sqlite) ListPics(album string, limit, page int) (ids []string) {
	query := "SELECT DISTINCT Hash FROM Albums WHERE Album = ? LIMIT ? OFFSET ?;"
	results, err := db.driver.Query(query, album, limit, limit*(page-1))
	if err != nil {
		log.Error(str.T("sql failed: can not get image id from album({})", album))
		return nil
	}
	for results.Next() {
		var id string
		err = results.Scan(&id)
		if err != nil {
			log.Error(str.T("sql failed: can not get image id from album({})", album))
			return nil
		}
		ids = append(ids, id)
	}
	log.Debug(str.T("find {} albums in database", len(ids)))
	return
}
