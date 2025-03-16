package pages

type ErrParams struct {
	Status  int
	Message string
}

type LinkResultParams struct {
	Shortlink string
	QRBase64  string
}

type ViewImageParams struct {
	FileName     string
	SizeMB       float32
	ViewUrl      string
	ImageUrl     string
	ThumbnailUrl string
}

type ViewFileParams struct {
	FileName        string
	FileSizeMB      float32
	FileViewUrl     string
	FileDownloadUrl string
}
