package image

import (
	"image/png"
	"os"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"github.com/EdlinOrg/prominentcolor"
)

// GetMainColor alert: 必须提供png格式图片, 返回6位HEX色值
func GetMainColor(path string) string {
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
	log.Debug(str.T("get the main color from image: {color}", color))
	return color
}
