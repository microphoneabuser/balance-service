package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/microphoneabuser/balance-service/models"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type balanceResponse struct {
	Balance float64 `json:"balance"`
}

type getTransactionsResponse struct {
	Data []models.TransactionForOutput `json:"data"`
}

func newErrorMessage(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
