package hub

import (
	"errors"

	"github.com/benbjohnson/immutable"
	"github.com/crosstyan/sh-over-ws/message"
	"github.com/crosstyan/sh-over-ws/utils"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Client interface {
	Uuid() uuid.UUID
	Conn() *websocket.Conn
}

type Hub struct {
	actuator   *immutable.Map[uuid.UUID, Actuator]
	controller *immutable.Map[uuid.UUID, Controller]
	// a temporary connection. Would be converted to `Actuator` or `Controller`.
	visitor *immutable.Map[uuid.UUID, Visitor]
}

// ugly... just
func (h *Hub) Get(uuid uuid.UUID) (Client, bool) {
	if v, ok := h.visitor.Get(uuid); ok {
		return v, ok
	}
	if v, ok := h.controller.Get(uuid); ok {
		return v, ok
	}
	if v, ok := h.actuator.Get(uuid); ok {
		return v, ok
	}
	return nil, false
}

func (h *Hub) NewVisitor(conn *websocket.Conn) uuid.UUID {
	v := Visitor{
		uuid: uuid.New(),
		conn: conn,
	}
	h.visitor = h.visitor.Set(v.Uuid(), v)
	return v.uuid
}

func (h *Hub) Remove(uuid uuid.UUID) {
	cl, ok := h.Get(uuid)
	if ok {
		switch c := cl.(type) {
		case Actuator:
			for _, cid := range c.Subscribers() {
				if controller, ok := h.controller.Get(cid); ok {
					controller.Unsubscribe(c.Uuid())
				}
			}
			h.actuator = h.actuator.Delete(uuid)
		case Controller:
			for _, aid := range c.Subscriptions() {
				if actuator, ok := h.actuator.Get(aid); ok {
					actuator.RemoveSubscriber(c.Uuid())
				}
			}
			h.controller = h.controller.Delete(uuid)
		case Visitor:
			h.visitor = h.visitor.Delete(uuid)
		default:
			// do nothing
		}
	}
}

func (h *Hub) Subscribe(req *message.ControlRequest) error {
	cid, err := uuid.FromBytes(req.ControllerId)
	if err != nil {
		return err
	}
	aid, err := uuid.FromBytes(req.ActuatorId)
	if err != nil {
		return err
	}
	actuator, ok := h.actuator.Get(aid)
	if !ok {
		return errors.New("actuator not found")
	}
	controller, ok := h.controller.Get(cid)
	if !ok {
		return errors.New("controller not found")
	}
	actuator.AddSubscriber(cid)
	controller.Subscribe(aid)
	return nil
}

func (h *Hub) FromVisitor(id uuid.UUID, handshake *message.Handshake) (uuid.UUID, error) {
	var err error
	empty := *new(uuid.UUID)
	switch hs := handshake.Handshake.(type) {
	case *message.Handshake_Actuator:
		id, err = h.toActuator(id, hs.Actuator)
	case *message.Handshake_Controller:
		id, err = h.toController(id, hs.Controller)
	}
	if err != nil {
		return empty, err
	}
	return id, nil
}

// TODO: too much duplicate code
func (h *Hub) toActuator(tempUuid uuid.UUID, handshake *message.ActuatorHandshake) (uuid.UUID, error) {
	realUuid, err := uuid.FromBytes(handshake.Uuid)
	empty := *new(uuid.UUID)
	if err != nil {
		return empty, err
	}
	name := handshake.Name
	v, ok := h.visitor.Get(tempUuid)
	if !ok {
		return empty, errors.New("not found")
	}
	h.actuator = h.actuator.Set(realUuid, v.ToActuator(realUuid, name))
	h.visitor = h.visitor.Delete(tempUuid)
	return realUuid, nil
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
	h.controller = h.controller.Set(realUuid, v.ToController(realUuid))
	h.visitor = h.visitor.Delete(tempUuid)
	return realUuid, nil
}

func NewHub() Hub {
	hasher := utils.UuidHasher{}
	return Hub{
		actuator:   immutable.NewMap[uuid.UUID, Actuator](&hasher),
		controller: immutable.NewMap[uuid.UUID, Controller](&hasher),
		visitor:    immutable.NewMap[uuid.UUID, Visitor](&hasher),
	}
}

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

func (c Controller) Subscribe(id uuid.UUID) {
	c.subscribing = append(c.subscribing, id)
}

func (c Controller) Unsubscribe(id uuid.UUID) {
	utils.DeleteIfOnce(
		c.subscribing,
		id,
		func(a, b uuid.UUID) bool {
			return a.ID() == b.ID()
		})
}

func (c Controller) Subscriptions() []uuid.UUID {
	return c.subscribing
}

type Actuator struct {
	name        string
	uuid        uuid.UUID
	conn        *websocket.Conn
	subscribers []uuid.UUID
}

func (a Actuator) Name() string {
	return a.name
}

func (a Actuator) Uuid() uuid.UUID {
	return a.uuid
}

func (a Actuator) Conn() *websocket.Conn {
	return a.conn
}

// same as `Controller.Subscribe`
func (a Actuator) AddSubscriber(id uuid.UUID) {
	a.subscribers = append(a.subscribers, id)
}

// same as `Controller.Unsubscribe`
func (a Actuator) RemoveSubscriber(id uuid.UUID) {
	utils.DeleteIfOnce(
		a.subscribers,
		id,
		func(a, b uuid.UUID) bool {
			return a.ID() == b.ID()
		})
}

func (a Actuator) Subscribers() []uuid.UUID {
	return a.subscribers
}
