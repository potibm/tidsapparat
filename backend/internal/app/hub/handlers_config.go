package hub

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/potibm/billedapparat/internal/app/config"
)

type AppConfigPublic struct {
	Version            string              `json:"version"`
	Environment        string              `json:"environment"`
	EnvironmentMessage string              `json:"environment_message"`
	Sentry             config.SentryConfig `json:"sentry"`
}

func (s *Server) handleGetPublicConfig(c *gin.Context) {
	pub := mapToPublicConfig(&s.cfg)

	c.JSON(http.StatusOK, pub)
}

func mapToPublicConfig(cfg *config.Config) AppConfigPublic {
	return AppConfigPublic{
		Version:            cfg.App.Version,
		Environment:        cfg.App.Environment,
		EnvironmentMessage: cfg.App.EnvironmentMessage,
		Sentry:             cfg.Sentry,
	}
}
