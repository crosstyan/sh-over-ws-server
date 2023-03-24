package hub

import (
	"errors"

	"github.com/crosstyan/sh-over-ws/message"
	"github.com/crosstyan/sh-over-ws/utils"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

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
func (a *Actuator) AddSubscriber(id uuid.UUID) {
	a.subscribers = append(a.subscribers, id)
}

// same as `Controller.Unsubscribe`
func (a *Actuator) RemoveSubscriber(id uuid.UUID) {
	utils.DeleteIfOnce(
		a.subscribers,
		id,
		func(a, b uuid.UUID) bool {
			return a.ID() == b.ID()
		})
}

func (a *Actuator) Subscribers() []uuid.UUID {
	return a.subscribers
}

func (h *Hub) toActuator(tempUuid uuid.UUID, handshake *message.ActuatorHandshake) (uuid.UUID, error) {
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
	name := handshake.Name
	h.actuator = h.actuator.Set(realUuid, v.ToActuator(realUuid, name))
	h.visitor = h.visitor.Delete(tempUuid)
	return realUuid, nil
}
