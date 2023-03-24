package hub

import (
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

// https://stackoverflow.com/questions/27775376/value-receiver-vs-pointer-receiver
// looks like you treat pointers like in C - for keeping side effects in
// operations. In Go there is another paradigm. The values are passed by value,
// and they keep pointers internally to share the state between copies.
type Visitor struct {
	uuid uuid.UUID
	conn *websocket.Conn
}

func (v Visitor) Uuid() uuid.UUID {
	return v.uuid
}

func (v Visitor) Conn() *websocket.Conn {
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

func (h *Hub) NewVisitor(conn *websocket.Conn) uuid.UUID {
	v := Visitor{
		uuid: uuid.New(),
		conn: conn,
	}
	h.visitor = h.visitor.Set(v.Uuid(), v)
	return v.uuid
}
