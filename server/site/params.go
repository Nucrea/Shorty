package site

type ErrParams struct {
	Status  int
	Message string
}

type LinkResultParams struct {
	Url string
}

type QRResultParams struct {
	ImageBase64 string
}

type ViewImageParams struct {
	FileName     string
	SizeMB       float32
	ViewUrl      string
	ImageUrl     string
	ThumbnailUrl string
}
