package main

// Should be based on https://github.com/sony/sonyflake
type generator struct {
	id int64
}

// NextID id generation
func (g *generator) NextID() int64 {
	g.id++
	return g.id
}
