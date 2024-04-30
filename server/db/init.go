package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	_ "github.com/mattn/go-sqlite3"
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

const (
	RpicDBPath = "./rpic/rpics.db"
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
	CreateTableStatistics = `CREATE TABLE IF NOT EXISTS RpicStat
(
    Request INTEGER,
    CacheCount INTEGER,
    ImageSpace INTEGER,
	CacheSpace INTEGER
);`
	InitTableStatistics  = `INSERT INTO RpicStat VALUES (0, 0, 0, 0);`
	CreateIndexImageDate = `CREATE INDEX IF NOT EXISTS idx_date ON Images(Date);`
	CreateIndexAlbum     = `CREATE INDEX IF NOT EXISTS idx_album ON Albums(Album);`
)

func Connect() *Sqlite {
	_ = os.MkdirAll(filepath.Dir(RpicDBPath), os.ModePerm)
	db, err := sql.Open("sqlite3", RpicDBPath)
	if err != nil {
		panic(str.Join("can not connect to database: ", RpicDBPath))
	}
	justexec(db, CreateTableImages)
	justexec(db, CreateTableAlbums)
	justexec(db, CreateTableImageData)
	justexec(db, CreateIndexImageDate)
	justexec(db, CreateIndexAlbum)
	justexec(db, CreateTableStatistics)
	row := db.QueryRow("SELECT COUNT(*) FROM RpicStat;")
	var count int
	var rpicStat *RpicStat
	err = row.Scan(&count)
	if err != nil {
		log.Error("sql error: can not get data from table(RpicStat)")
	} else {
		if count == 0 {
			justexec(db, InitTableStatistics)
		}
		// info: read table RpicStat
		var requestCount, cacheCount int
		var imageSpace, cacheSpace int64
		row = db.QueryRow("SELECT Request, CacheCount, ImageSpace, CacheSpace FROM RpicStat;")
		err = row.Scan(&requestCount, &cacheCount, &imageSpace, &cacheSpace)
		if err != nil {
			log.Error("sql error: can not get data from table(RpicStat)")
		}
		rpicStat = &RpicStat{
			requestCount: requestCount,
			cacheCount:   cacheCount,
			imageSpace:   imageSpace,
			cacheSpace:   cacheSpace,
			runner:       false,
			RequestAdder: make(chan struct{}, 10),
			CacheAdder: make(
				chan struct {
					contentSize int64
					isCache     bool
				}, 10,
			),
		}
	}
	sqliteDB := &Sqlite{
		driver:   db,
		RpicStat: rpicStat,
	}
	sqliteDB.RunStat()
	return sqliteDB
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
