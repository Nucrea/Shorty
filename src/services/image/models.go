package image

type ImageMetadataDTO struct {
	Id          string
	Name        string
	OriginalId  string
	ThumbnailId string
}

type ImageMetadataExDTO struct {
	Id          string
	Name        string
	Size        int
	Hash        string
	OriginalId  string
	ThumbnailId string
}
