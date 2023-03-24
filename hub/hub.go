package hub

import (
	"context"
	"errors"

	"github.com/benbjohnson/immutable"
	"github.com/crosstyan/sh-over-ws/log"
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

// https://stackoverflow.com/questions/27775376/value-receiver-vs-pointer-receiver
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

func (h *Hub) HandleStdPayload(ctx context.Context, payload *message.StdPayload) error {
	broadcast := func(ctx context.Context, a *Actuator, b []byte) {
		for _, s := range a.subscribers {
			sub, ok := h.Get(s)
			if ok {
				go func() {
					err := sub.Conn().Write(ctx, websocket.MessageBinary, b)
					if err != nil {
						log.Sugar().Errorw("WriteError", "type", "STDOUT", "error", err)
					}
				}()
			}
		}
	}

	aid, err := uuid.FromBytes(payload.Uuid)
	if err != nil {
		return err
	}
	actuator, ok := h.actuator.Get(aid)
	if !ok {
		return errors.New("actuator not found")
	}

	switch p := payload.Payload.(type) {
	case *message.StdPayload_Stdout:
		broadcast(ctx, &actuator, p.Stdout.Data)
		return nil
	case *message.StdPayload_Stdin:
		err = actuator.Conn().Write(ctx, websocket.MessageBinary, p.Stdin.Data)
		return err
	}
	return nil
}

func (h *Hub) RequestSub(msg *message.ControlRequest) error {
	cid, err := uuid.FromBytes(msg.ControllerId)
	if err != nil {
		return err
	}
	aid, err := uuid.FromBytes(msg.ActuatorId)
	if err != nil {
		return err
	}
	actuator, ok := h.actuator.Get(aid)
	if !ok {
		return errors.New("No actuator")
	}
	controller, ok := h.controller.Get(cid)
	if !ok {
		return errors.New("No controller")
	}
	actuator.AddSubscriber(cid)
	controller.Subscribe(aid)
	return nil
}

func NewHub() Hub {
	hasher := utils.UuidHasher{}
	return Hub{
		actuator:   immutable.NewMap[uuid.UUID, Actuator](&hasher),
		controller: immutable.NewMap[uuid.UUID, Controller](&hasher),
		visitor:    immutable.NewMap[uuid.UUID, Visitor](&hasher),
	}
}
