package image

type ImageInfoDTO struct {
	Id          string
	ShortId     string
	Size        int
	Hash        string
	Name        string
	ImageId     string
	ThumbnailId string
}

type ImageDTO struct {
	Id    string
	Name  string
	Size  int
	Bytes []byte
}
