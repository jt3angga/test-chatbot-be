package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jt3angga/test-chatbot-be/internal/client"
)

type ChatHandler struct {
    GroqClient *client.GroqClient
}

func NewChatHandler(gc *client.GroqClient) *ChatHandler {
    return &ChatHandler{GroqClient: gc}
}

func (ch *ChatHandler) ChatStream(c *gin.Context) {
    var request struct {
        Message string `json:"message"`
    }

    if err := c.BindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    stream, err := ch.GroqClient.StreamResponse(c.Request.Context(), request.Message)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer stream.Close()

    c.Stream(func(w io.Writer) bool {
        buf := make([]byte, 1024)
        for {
            n, err := stream.Read(buf)
            if err != nil {
                return false
            }
            if n > 0 {
                c.Writer.Write(buf[:n])
                c.Writer.Flush()
            }
        }
    })
}