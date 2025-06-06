package common

type ResourcePathable interface {
	CollectionPath() string
	ResourcePath(id string) string
}
