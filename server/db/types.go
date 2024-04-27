package db

type Image struct {
	Hash  string
	Scale string
	Date  int
}

type ImageData struct {
	Path   string
	Hash   string
	Size   string
	Width  int
	Height int
	Format string
}

type RpicRequest struct {
	Album  string
	Scale  string
	HasRid bool
	Rid    int // TODO: make sure rid > 0
}
