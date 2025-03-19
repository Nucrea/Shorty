package assets

type AssetMetadataDTO struct {
	Id   string
	Size int
	Hash string
}

type AssetDTO struct {
	Id    string
	Size  int
	Hash  string
	Bytes []byte
}

type AssetStatus string

const (
	AssetPending AssetStatus = "pending"
	AssetCreated AssetStatus = "created"
	AssetDeleted AssetStatus = "deleted"
)
