package wnet

type ConnectionID struct {
	id uint32
}

func NewConnectionID() *ConnectionID {
	return &ConnectionID{id: 0}
}

func (c *ConnectionID) Get() uint32 {
	return c.id
}

func (c *ConnectionID) Inc() uint32 {
	c.id++
	return c.id
}
