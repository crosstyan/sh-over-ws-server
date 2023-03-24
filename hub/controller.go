package hub

import (
	"errors"

	"github.com/crosstyan/sh-over-ws/message"
	"github.com/crosstyan/sh-over-ws/utils"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Controller struct {
	uuid        uuid.UUID
	conn        *websocket.Conn
	subscribing []uuid.UUID
}

func (c Controller) Uuid() uuid.UUID {
	return c.uuid
}

func (c Controller) Conn() *websocket.Conn {
	return c.conn
}

func (c *Controller) Subscribe(id uuid.UUID) {
	c.subscribing = append(c.subscribing, id)
}

func (c *Controller) Unsubscribe(id uuid.UUID) {
	utils.DeleteIfOnce(
		c.subscribing,
		id,
		func(a, b uuid.UUID) bool {
			return a.ID() == b.ID()
		})
}

func (c *Controller) Subscriptions() []uuid.UUID {
	return c.subscribing
}

func (h *Hub) toController(tempUuid uuid.UUID, handshake *message.ControllerHandshake) (uuid.UUID, error) {
	realUuid, err := uuid.FromBytes(handshake.Uuid)
	empty := *new(uuid.UUID)
	if err != nil {
		return empty, err
	}
	v, ok := h.visitor.Get(tempUuid)
	if !ok {
		return empty, errors.New("not found")
	}
	_, ok = h.Get(realUuid)
	if ok {
		return empty, errors.New("uuid already exists")
	}
	h.controller = h.controller.Set(realUuid, v.ToController(realUuid))
	h.visitor = h.visitor.Delete(tempUuid)
	return realUuid, nil
}
