package db

func (db *Sqlite) RemoveAll(hash string) (ok bool) {
	query := []string{
		`DELETE FROM ImageData WHERE Hash = ?;`, `DELETE FROM Albums WHERE Hash = ?;`,
		`DELETE FROM Images WHERE Hash = ?;`,
	}
	for i := range 3 {
		_, err := db.driver.Query(query[i], hash)
		if err != nil {
			return false
		}
	}
	return true
}

func (db *Sqlite) RemoveFromAlbum(hash, album string) (ok bool) {
	query := "DELETE FROM Albums WHERE Hash = ? AND Album = ?;"
	_, err := db.driver.Query(query, hash, album)
	if err != nil {
		return false
	}
	return true
}
