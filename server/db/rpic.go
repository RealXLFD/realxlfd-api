package db

import (
	"database/sql"
	"errors"
	"regexp"
	"strconv"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

func (db *Sqlite) Rpic(req *RpicRequest) (hash string, main string, ok bool) {
	template := "SELECT Images.Hash, Images.Main FROM Images JOIN Albums ON Images.Hash = Albums.Hash {} ORDER BY {} LIMIT 1 {};"
	var dArgOrder, dArgOffset string
	var result *sql.Row
	if !req.HasRid {
		dArgOrder = "RANDOM()"
		dArgOffset = ""
	} else {
		dArgOrder = "Images.Date"
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
	scale := ParseScale(req.Scale)
	switch {
	case req.Album != "" && scale != "":
		query := str.T(
			template,
			str.Join("WHERE Albums.Album = ? AND", scale), dArgOrder,
			dArgOffset,
		)
		result = db.driver.QueryRow(query, req.Album, req.Scale)
	case req.Album != "":
		query := str.T(
			template, "WHERE Albums.Album = ?", dArgOrder,
			dArgOffset,
		)
		result = db.driver.QueryRow(query, req.Album)
	case scale != "":
		query := str.T(
			template, str.Join("WHERE", scale), dArgOrder,
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

func ParseScale(scale string) (query string) {
	switch scale {
	case "pc":
		return " Images.Scale > 1"
	case "mobile":
		return " Images.Scale < 1"
	default:
		re := regexp.MustCompile(`([<>])(-?\d+(\.\d+)?)`)
		matches := re.FindStringSubmatch(scale)
		if matches == nil {
			return ""
		}
		operation := matches[1]
		scaleNum, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return ""
		}
		return str.T(
			" Images.Scale {sign} {number}",
			operation,
			strconv.FormatFloat(scaleNum, 'f', 2, 64),
		)
	}
}
