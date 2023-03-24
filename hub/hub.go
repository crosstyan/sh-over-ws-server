package hub

import (
	"errors"

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
	// a temporary connection. Would be converted to `Actuator` or `Controller`.
	Visitor *immutable.Map[uuid.UUID, Visitor]
}

// TODO: too much duplicate code
func (h *Hub) ToActuator(tempUuid uuid.UUID, realUuid uuid.UUID, name string) error {
	v, ok := h.Visitor.Get(tempUuid)
	if !ok {
		return errors.New("not found")
	}
	h.Actuator = h.Actuator.Set(realUuid, v.ToActuator(realUuid, name))
	h.Visitor = h.Visitor.Delete(tempUuid)
	return nil
}

func (h *Hub) ToController(tempUuid uuid.UUID, realUuid uuid.UUID) error {
	v, ok := h.Visitor.Get(tempUuid)
	if !ok {
		return errors.New("not found")
	}
	h.Controller = h.Controller.Set(realUuid, v.ToController(realUuid))
	h.Visitor = h.Visitor.Delete(tempUuid)
	return nil
}

func NewHub() Hub {
	hasher := UuidHasher{}
	return Hub{
		Actuator:   immutable.NewMap[uuid.UUID, Actuator](&hasher),
		Controller: immutable.NewMap[uuid.UUID, Controller](&hasher),
		Visitor:    immutable.NewMap[uuid.UUID, Visitor](&hasher),
	}
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
