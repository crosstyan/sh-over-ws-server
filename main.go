package main

import (
	"strings"

	"github.com/crosstyan/sh-over-ws/hub"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

// https://pkg.go.dev/golang.org/x/net/websocket
// official recommendation is
// https://pkg.go.dev/nhooyr.io/websocket
func main() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	h := hub.NewHub()
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			sugar.Errorw("WebsocketAccept", "error", err, "from", c.Request.RemoteAddr)
		}
		// no idea what go context is
		// https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/
		v := hub.NewVisitor(conn)
		h.Visitor = h.Visitor.Set(v.Uuid(), v)
		go func(done <-chan struct{}) {
			<-done
			h.Visitor = h.Visitor.Delete(v.Uuid())
			sugar.Infow("VisitorDeleted", "from", c.Request.RemoteAddr)
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
				// Go declaration syntax says nothing about stack or heap.
				l, err := reader.Read(buffer)
				if err != nil {
					sugar.Errorw("ReaderRead", "error", err, "from", c.Request.RemoteAddr)
				}
				m := string(buffer[:l])
				m = strings.Trim(m, " \t")
				if m != "" {
					sugar.Warnw("MessageText", "content", m, "from", c.Request.RemoteAddr)
					// Echo. Nothing fancy.
					_ = conn.Write(c, websocket.MessageText, []byte(m))
				}
			case websocket.MessageBinary:
				// do nothing
			}
		}
	})
	r.Run()
}
