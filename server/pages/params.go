package pages

type ErrParams struct {
	Status  int
	Message string
}

type LinkFormParams struct {
	Id            string
	CaptchaBase64 string
}

type ImageFormParams struct {
	Id            string
	CaptchaBase64 string
}

type LinkResultParams struct {
	Shortlink string
	QRBase64  string
}

type ImageViewParams struct {
	FileName     string
	SizeMB       float32
	ViewUrl      string
	ImageUrl     string
	ThumbnailUrl string
}

type FileViewParams struct {
	FileName        string
	FileSizeMB      float32
	FileViewUrl     string
	FileDownloadUrl string
	CaptchaId       string
	CaptchaBase64   string
}

type FileDownloadParams struct {
	FileRawUrl string
}
