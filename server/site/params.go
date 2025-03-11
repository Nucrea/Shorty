package site

type ErrParams struct {
	Code    int
	Message string
}

type LinkResultParams struct {
	Url string
}

type QRResultParams struct {
	ImageBase64 string
}

type ViewImageParams struct {
	FileName string
	Size     int
	Url      string
}
