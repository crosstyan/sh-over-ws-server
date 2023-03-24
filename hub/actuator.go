package hub

import (
	"context"
	"errors"
	"time"

	"github.com/crosstyan/sh-over-ws/log"
	"github.com/crosstyan/sh-over-ws/message"
	"github.com/crosstyan/sh-over-ws/utils"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

type Actuator struct {
	name        string
	uuid        uuid.UUID
	conn        *websocket.Conn
	subscribers []uuid.UUID
	events      <-chan message.ControlState
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
	if len(a.Subscribers()) <= 0 {
		ctx := context.Background()
		timeout, _ := time.ParseDuration("500ms")
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		err := a.SendState(ctx, message.ControlState_BIND)
		log.Sugar().Errorw("actuator.SendBind", "error", err)
	}
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

// Side effect: Websocket
func (a *Actuator) SendState(ctx context.Context, state message.ControlState) error {
	msg := new(message.ServerMsg)
	ctl := new(message.ServerMsg_ActuatorControl)
	msg.Payload = ctl
	ctl.ActuatorControl.State = state
	b, err := a.Uuid().MarshalBinary()
	if err != nil {
		return err
	}

	ctl.ActuatorControl.Uuid = b
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	err = a.Conn().Write(ctx, websocket.MessageBinary, data)
	if err != nil {
		return err
	}
	return nil
}
