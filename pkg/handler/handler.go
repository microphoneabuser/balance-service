package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/microphoneabuser/balance-service/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) SetupRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/balance", h.getBalance)
	router.POST("/accrual", h.postAccrual)
	router.POST("/debiting", h.postDebiting)
	router.POST("/transfer", h.postTransfer)
	router.GET("/transactions", h.getTransactions)

	return router
}
