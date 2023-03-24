package main

import (
	"strings"
	"time"

	"github.com/crosstyan/sh-over-ws/hub"
	"github.com/crosstyan/sh-over-ws/log"
	"github.com/crosstyan/sh-over-ws/message"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

func main() {
	sugar := log.Sugar()
	logger := log.Logger()
	h := hub.NewHub()
	r := gin.New()
	r.Use(ginZap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginZap.RecoveryWithZap(logger, true))
	r.GET("/ws", func(c *gin.Context) {
		conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			sugar.Warnw("WebsocketAccept", "error", err, "from", c.Request.RemoteAddr)
		}
		// NOTE: mutable
		id := h.NewVisitor(conn)
		go func(done <-chan struct{}) {
			<-done
			h.Remove(id)
			sugar.Infow("ClientLeave", "from", c.Request.RemoteAddr)
		}(c.Request.Context().Done())
		for {
			t, reader, err := conn.Reader(c.Request.Context())
			if err != nil {
				conn.Close(websocket.StatusInternalError, "Failed to get reader")
				return
			}
			buffer := make([]byte, 1024, 10240)
			switch t {
			case websocket.MessageText:
				l, err := reader.Read(buffer)
				if err != nil {
					sugar.Errorw("ReaderRead", "error", err, "from", c.Request.RemoteAddr)
				}
				m := string(buffer[:l])
				m = strings.Trim(m, " \t")
				if m != "" {
					sugar.Infow("MessageText", "content", m, "from", c.Request.RemoteAddr)
					// Echo. Nothing fancy.
					_ = conn.Write(c, websocket.MessageText, []byte(m))
				}
			case websocket.MessageBinary:
				payload := &message.ClientMsg{}
				err := proto.Unmarshal(buffer, payload)
				if err != nil {
					sugar.Errorw("ProtoUnmarshal", "error", err, "from", c.Request.RemoteAddr)
				}
				switch p := payload.Payload.(type) {
				case *message.ClientMsg_Handshake:
					uid, err := h.FromVisitor(id, p.Handshake)
					sugar.Infow("FromVisitor", "uid", uid, "from", c.Request.RemoteAddr)
					if err != nil {
						sugar.Errorw("FromVisitor", "error", err, "from", c.Request.RemoteAddr)
					}
					id = uid
				case *message.ClientMsg_ControlRequest:
					err = h.RequestSub(p.ControlRequest)
					sugar.Infow("RequestSub", "from", c.Request.RemoteAddr)
					if err != nil {
						sugar.Errorw("RequestSub", "error", err, "from", c.Request.RemoteAddr)
					}
				case *message.ClientMsg_StdPayload:
					err = h.HandleStdPayload(c, p.StdPayload)
					sugar.Infow("HandleStdPayload", "from", c.Request.RemoteAddr)
					if err != nil {
						sugar.Errorw("HandleStdPayload", "error", err, "from", c.Request.RemoteAddr)
					}
				}
			}
		}
	})
	r.Run()
}
