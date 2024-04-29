package image

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
)

func MD5(src string) string {
	file, err := os.Open(src)
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
	if _, err = io.Copy(hash, file); err != nil {
		log.Error(str.T("error occurred while generate md5 from: {src}", src))
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}
