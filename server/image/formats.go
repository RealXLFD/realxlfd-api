package image

var (
	Formats = map[string]string{
		"webp": "image/webp",
		"jpeg": "image/jpeg",
		"jpg":  "image/jpeg",
		"gif":  "image/gif",
		"png":  "image/png",
		"avif": "image/avif",
		"tiff": "image/tiff",
	}
	Quality = map[string]int{
		"lossless": 5,
		"high":     4,
		"default":  3,
		"normal":   2,
		"mid":      1,
		"low":      0,
		"5":        5,
		"4":        4,
		"3":        3,
		"2":        2,
		"1":        1,
		"0":        0,
	}
	QualityArr = []string{
		"[Q=60]", "[Q=70]", "[Q=80]", "[Q=90]", "[Q=95]", "",
	}
	Sizes = map[string]string{
		"raw":   "",
		"4k":    "3840",
		"2k":    "2560",
		"1080p": "1080",
		"720p":  "720",
	}
)
