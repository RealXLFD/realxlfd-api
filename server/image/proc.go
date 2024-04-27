package image

import (
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

type Info struct {
	file      string
	Prominent string
	Scale     string
	Format    string
}

func GetInfo(files ...string) (infos []Info, ok bool) {

}

func parseInfo(file, vipsoutput string) Info {

}
