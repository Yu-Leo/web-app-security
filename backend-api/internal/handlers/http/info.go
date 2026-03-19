package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInfo(c *gin.Context) {
	c.Status(http.StatusOK)
}
