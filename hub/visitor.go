package hub

import (
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

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

// the original visitor should be removed from the hub
func (v *Visitor) ToActuator(uuid uuid.UUID, name string) Actuator {
	return Actuator{
		uuid: uuid,
		conn: v.conn,
		name: name,
	}
}

func (v *Visitor) ToController(uuid uuid.UUID) Controller {
	return Controller{
		uuid: uuid,
		conn: v.conn,
	}
}

func NewVisitor(conn *websocket.Conn) Visitor {
	return Visitor{
		uuid: uuid.New(),
		conn: conn,
	}
}
