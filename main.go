package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
)

// https://pkg.go.dev/golang.org/x/net/websocket
// official recommendation is
// https://pkg.go.dev/nhooyr.io/websocket
func main() {
	r := gin.Default()
	fmt.Println("Hello, World!")
	r.GET("/ws", func(c *gin.Context) {
		conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{})
		if err != nil {
			// copilot did this
			conn.Close(websocket.StatusInternalError, "Failed to accept websocket connection")
		}
	})
}
