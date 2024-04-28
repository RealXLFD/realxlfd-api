package image

import (
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"rpics-docker/server"
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

type Info map[string]metadata
type metadata struct {
	Prominent string
	Scale     float64
}

// GetInfo alert 必须提供PNG格式图片
func GetInfo(files ...string) (info Info) {
	info = make(Info)
	files = append(files, str.T("--vips-concurrency={}", runtime.NumCPU()))
	cmd := exec.Command("vipsheader", files...)
	out, err := cmd.Output()
	if err != nil {
		log.Error("error occurred during getting image header")
	}
	re := regexp.MustCompile(`^(.*): (\d+)x(\d+)`)
	lines := strings.Split(string(out), "\r\n")
	if len(lines) == 1 {
		log.Warn("invalid image: {file}", files[0])
		return nil
	}
	task := sync.WaitGroup{}
	task.Add(len(lines) - 1)
	for i := range len(lines) - 1 {
		matches := re.FindStringSubmatch(lines[i])
		if matches == nil {
			log.Error("can not parse image header info")
			return nil
		}
		filename := matches[1]
		width, _ := strconv.Atoi(matches[2])
		height, _ := strconv.Atoi(matches[3])
		scale := float64(width) / float64(height)
		server.ThreadPool.Push(
			func() {
				color := GetMainColor(filename)
				if color != "" {
					info[filename] = metadata{color, scale}
				}
				task.Done()
			},
		)
	}
	task.Wait()
	return
}

func Convert(src, dst, format, size, quality string) {

}
