package files

type FileMetadataDTO struct {
	Id     string
	FileId string
	Name   string
}

type FileMetadataExDTO struct {
	Id     string
	FileId string
	Name   string
	Size   int
	Hash   string
}
