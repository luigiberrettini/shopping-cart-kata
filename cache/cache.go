package cache

// Etagger represents a type that has etag
type Etagger interface {
	GetEtag() string
	ComputeEtag()
}

// Cache represents the cache for cart resources
type Cache interface {
	GetByEtagWithID(etag string, wid string) (Etagger, bool)
	AddOrReplace(wid string, e Etagger)
	Remove(wid string)
}
