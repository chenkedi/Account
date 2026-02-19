package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StatsHandler struct {
	logger *zap.Logger
}

func NewStatsHandler(logger *zap.Logger) *StatsHandler {
	return &StatsHandler{logger: logger}
}
