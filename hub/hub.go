package hub

import (
	"bytes"

	"github.com/benbjohnson/immutable"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Client interface {
	Uuid() uuid.UUID
	Conn() *websocket.Conn
}

type Hub struct {
	Actuator   *immutable.Map[uuid.UUID, Actuator]
	Controller *immutable.Map[uuid.UUID, Controller]
	// would remove this in the future
	Visitor *immutable.Map[uuid.UUID, Visitor]
}

func NewHub() Hub {
	hasher := UuidHasher{}
	return Hub{
		Actuator:   immutable.NewMap[uuid.UUID, Actuator](&hasher),
		Controller: immutable.NewMap[uuid.UUID, Controller](&hasher),
		Visitor:    immutable.NewMap[uuid.UUID, Visitor](&hasher),
	}
}

type Visitor struct {
	uuid uuid.UUID
	conn *websocket.Conn
}

func (v *Visitor) Uuid() uuid.UUID {
	return v.uuid
}

func (v *Visitor) Conn() *websocket.Conn {
	return v.conn
}

func NewVisitor(conn *websocket.Conn) Visitor {
	return Visitor{
		uuid: uuid.New(),
		conn: conn,
	}
}

type UuidHasher struct{}

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 4122, which is way bigger than uint32.
func (h *UuidHasher) Hash(key uuid.UUID) uint32 {
	return key.ID()
}

func (h *UuidHasher) Equal(a, b uuid.UUID) bool {
	return bytes.Equal(a[:], b[:])
}

type Controller struct {
	uuid uuid.UUID
	conn *websocket.Conn
}

func (c *Controller) Uuid() uuid.UUID {
	return c.uuid
}

func (c *Controller) Conn() *websocket.Conn {
	return c.conn
}

type Actuator struct {
	name string
	uuid uuid.UUID
	conn *websocket.Conn
}

func (a *Actuator) Name() string {
	return a.name
}

func (a *Actuator) Uuid() uuid.UUID {
	return a.uuid
}

func (a *Actuator) Conn() *websocket.Conn {
	return a.conn
}
