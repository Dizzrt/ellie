package ginx

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func TestDecode(t *testing.T) {
	r := gin.Default()
	r.POST("/user/:id", func(ctx *gin.Context) {
		var req UserRequest
		if err := DecodeRequest(ctx, &req); err != nil {
			t.Errorf("DecodeRequest error: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Decoded Request: %+v\n", req)
		ctx.JSON(http.StatusOK, gin.H{
			"user": req,
		})
	})

	r.Run()
}
