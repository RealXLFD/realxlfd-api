package image

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"git.realxlfd.cc/RealXLFD/golib/utils/str"
	"rpics-docker/serverlog"
)

var (
	log = serverlog.Log
)

type Info map[string]Metadata
type Metadata struct {
	Prominent string
	Scale     float64
	Width     int
	Height    int
}

// GetInfo alert: 必须提供PNG格式图片, 确保输入路径为绝对路径
func GetInfo(file string) (info Metadata) {
	cmd := exec.Command("vipsheader", file)
	out, err := cmd.Output()
	if err != nil {
		log.Error("error occurred during getting image header")
	}
	re := regexp.MustCompile(`^(.*): (\d+)x(\d+)`)
	matches := re.FindStringSubmatch(string(out))
	if matches == nil {
		log.Error("can not parse image header infos")
		return
	}
	width, _ := strconv.Atoi(matches[2])
	height, _ := strconv.Atoi(matches[3])
	scale := float64(width) / float64(height)
	color := GetMainColor(file)
	info = Metadata{color, scale, width, height}
	return
}

// GetAllInfos alert: 必须提供PNG格式图片, 确保输入路径为绝对路径
func GetAllInfos(files ...string) (infos Info) {
	infos = make(Info)
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
	for i := range len(lines) - 1 {
		matches := re.FindStringSubmatch(lines[i])
		if matches == nil {
			log.Error("can not parse image header infos")
			return nil
		}
		filename := matches[1]
		width, _ := strconv.Atoi(matches[2])
		height, _ := strconv.Atoi(matches[3])
		scale := float64(width) / float64(height)
		color := GetMainColor(filename)
		if color != "" {
			infos[filename] = Metadata{color, scale, width, height}
		}
	}
	return
}

// Convert alert: 确保输入路径为绝对路径
func Convert(src, dst, size string, quality int) (ok bool) {
	// info: vips thumbnail <src> <dst>[Q=?/lossless] <Width> --Height <pixels>
	// info: vips copy <src> <dst>[Q=?/lossless]
	var vipsQuality, vipsPixels string
	var cmd *exec.Cmd
	vipsQuality = QualityArr[quality]
	switch size {
	case "raw":
		cmd = exec.Command("vips", "copy", src, str.Join(dst, vipsQuality))
	default:
		vipsPixels = Sizes[size]
		cmd = exec.Command(
			"vips",
			"thumbnail",
			src,
			str.Join(dst, vipsQuality),
			vipsPixels,
			"--height",
			vipsPixels,
			"--size",
			"down",
		)
	}
	log.Debug(str.T("execute: {cmd}", cmd.String()))
	err := cmd.Run()
	if err != nil {
		log.Error(str.T("vips error: can not convert {file}", filepath.Base(src)))
		return false
	}
	log.Debug(str.T("convert {src} to {dst}", src, dst))
	return true
}
