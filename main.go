package main

import (
	"image/png"
	"os"

	"git.realxlfd.cc/RealXLFD/golib/cli/logger"
	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/EdlinOrg/prominentcolor"
	_ "github.com/mattn/go-sqlite3"
)

var (
	log = logger.New()
)

func main() {
	println(str.T("{hello},123", "nihao"))
}

func test(path string) string {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error(str.T("can not open image file: {path}", path))
		return ""
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	image, err := png.Decode(file)
	if err != nil {
		log.Error(str.T("invalid image: {path}", path))
		return ""
	}
	mainColors, err := prominentcolor.KmeansWithArgs(1, image)
	if err != nil {
		log.Error(str.T("can not get main color from image: {path}", path))
		return ""
	}
	color := mainColors[0].AsString()
	return color
}
