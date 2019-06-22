package cache

import "sync"

// InMemCache is an in memory cache
type InMemCache struct {
	sync.RWMutex
	entriesByID    map[string]Etagger
	entryIdsByEtag map[string]string
}

// NewCache creates a new cache
func NewCache() Cache {
	c := new(InMemCache)
	c.entriesByID = make(map[string]Etagger)
	c.entryIdsByEtag = make(map[string]string)
	return c
}

// GetByEtagWithID get an entry by etags if its id matches the one provided
func (c *InMemCache) GetByEtagWithID(etag string, wid string) (Etagger, bool) {
	c.RLock()
	defer c.RUnlock()
	id, ok := c.entryIdsByEtag[etag]
	var e Etagger
	if !ok || id != wid {
		return e, false
	}
	e, ok = c.entriesByID[id]
	if !ok {
		return e, false
	}
	return e, true
}

// AddOrReplace adds or replaces an entry
func (c *InMemCache) AddOrReplace(wid string, e Etagger) {
	c.Lock()
	defer c.Unlock()
	remove(c, wid)
	c.entriesByID[wid] = e
	c.entryIdsByEtag[e.GetEtag()] = wid
}

// Remove removes an entry
func (c *InMemCache) Remove(wid string) {
	c.Lock()
	defer c.Unlock()
	remove(c, wid)
}

func remove(c *InMemCache, wid string) {
	e, ok := c.entriesByID[wid]
	if ok {
		delete(c.entryIdsByEtag, e.GetEtag())
	}
	delete(c.entriesByID, wid)
}
