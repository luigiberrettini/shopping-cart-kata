package cache

import (
	"fmt"
	"testing"
)

func TestEmptyCacheMiss(t *testing.T) {
	const (
		id    = "myID"
		value = "myValue"
	)
	e := etagger{ID: id, Value: value}
	e.ComputeEtag()
	c := NewCache()
	h, ok := c.GetByEtagWithID(e.etag, id)
	if ok {
		t.Errorf("Cache hit %v on not stored entry %v", h, e)
	}
}

func TestHitAfterAdd(t *testing.T) {
	const (
		id    = "myID"
		value = "myValue"
	)
	e := &etagger{ID: id, Value: value}
	e.ComputeEtag()
	c := NewCache()
	c.AddOrReplace(id, e)
	_, ok := c.GetByEtagWithID(e.etag, id)
	if !ok {
		t.Errorf("Cache miss on stored entry %v", e)
	}
}

func TestHitAfterReplace(t *testing.T) {
	const (
		id     = "myID"
		value1 = "myValue"
		value2 = "myValue"
	)
	e1 := &etagger{ID: id, Value: value1}
	e1.ComputeEtag()
	e2 := &etagger{ID: id, Value: value2}
	e2.ComputeEtag()
	c := NewCache()
	c.AddOrReplace(id, e1)
	c.AddOrReplace(id, e2)
	_, ok := c.GetByEtagWithID(e2.etag, id)
	if !ok {
		t.Errorf("Cache miss on replaced entry %v", e2)
	}
}

func TestRemove(t *testing.T) {
	const (
		id    = "myID"
		value = "myValue"
	)
	e := &etagger{ID: id, Value: value}
	e.ComputeEtag()
	c := NewCache()
	c.AddOrReplace(id, e)
	c.Remove(id)
	_, ok := c.GetByEtagWithID(e.etag, id)
	if ok {
		t.Errorf("Cache hit on remove entry %v", e)
	}
}

type etagger struct {
	ID    string
	Value string
	etag  string
}

func (e *etagger) GetEtag() string {
	return e.etag
}

func (e *etagger) ComputeEtag() {
	data := fmt.Sprintf("%s-%s", e.ID, e.Value)
	e.etag = fmt.Sprintf(`W/"%s-%d-%08X"`, "etagger", len(data), data)
}
