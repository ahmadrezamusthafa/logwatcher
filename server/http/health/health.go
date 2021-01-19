package health

import (
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Config config.Config `inject:"config"`
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}
