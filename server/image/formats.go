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
	Quality = map[string]string{
		"lossless": "lossless",
		"high":     "Q=95",
		"default":  "Q=90",
		"normal":   "Q=80",
		"mid":      "Q=70",
		"low":      "Q=60",
	}
	QualityArr = []string{
		"Q=60", "Q=70", "Q=80", "Q=90", "Q=95", "lossless",
	}
)
