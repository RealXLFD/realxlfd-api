package db

import (
	"database/sql"
	"strconv"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	_ "github.com/mattn/go-sqlite3"
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

const (
	Path = "./rpics.db"
)

const (
	CreateTableImages = `
CREATE TABLE IF NOT EXISTS Images
(
    Hash  VARCHAR(32) PRIMARY KEY,
    Main  VARCHAR(6),
    Scale REAL,
    Date  INTEGER
);`
	CreateTableAlbums = `
CREATE TABLE IF NOT EXISTS Albums
(
    Hash VARCHAR(32),
    Album VARCHAR(24),
    UNIQUE(Hash,Album),
	FOREIGN KEY (Hash) REFERENCES Images(Hash)
);
`
	CreateTableImageData = `
CREATE TABLE IF NOT EXISTS ImageData
(
    Path   TEXT PRIMARY KEY,
    Hash   VARCHAR(32) NOT NULL ,
    Size   VARCHAR(5) NOT NULL ,
    Quality INTEGER,
    Format VARCHAR(4),
    ContentSize INTEGER,
    FOREIGN KEY (Hash) REFERENCES Images (Hash)
);`
	CreateIndexImageDate = `CREATE INDEX IF NOT EXISTS idx_date ON Images(Date);`
	CreateIndexAlbum     = `CREATE INDEX IF NOT EXISTS idx_album ON Albums(Album);`
)

func Connect() *Sqlite {
	db, err := sql.Open("sqlite3", Path)
	if err != nil {
		panic(str.Join("can not connect to database: ", Path))
	}
	justexec(db, CreateTableImages)
	justexec(db, CreateTableAlbums)
	justexec(db, CreateTableImageData)
	justexec(db, CreateIndexImageDate)
	justexec(db, CreateIndexAlbum)
	return &Sqlite{
		driver: db,
	}
}

func justexec(db *sql.DB, state string) {
	stmt, err := db.Prepare(state)
	if err != nil {
		panic(err)
	}
	result, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
	affectRow, _ := result.RowsAffected()
	log.Debug(
		`execute `, state, ` ,`, strconv.Itoa(int(affectRow)), ` rows affected`,
	)
}
