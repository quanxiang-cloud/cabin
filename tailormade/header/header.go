package header

type CTX interface {
	GetHeader(key string) string
}
