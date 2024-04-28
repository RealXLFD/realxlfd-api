package image

import (
	"crypto/md5"
	"io"
	"os"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

func MD5(src string) string {
	file, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error(str.T("can not open file: {src}", src))
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("err occurred while closing file: {file}")
		}
	}(file)
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Error(str.T("error occurred while generate md5 from: {src}", src))
		return ""
	}
	return str.F("%x", md5.Sum(nil))
}
