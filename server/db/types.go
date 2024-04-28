package db

type Image struct {
	Hash  string
	Main  string
	Scale float64
	Date  int
}

type ImageData struct {
	Path    string
	Hash    string
	Size    string
	Quality int
	Format  string
}

type RpicRequest struct {
	Album  string
	Scale  string
	HasRid bool
	Rid    int // TODO: make sure rid > 0
}
