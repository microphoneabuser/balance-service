package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microphoneabuser/balance-service/models"
	"github.com/microphoneabuser/balance-service/rabbitmq"
)

func (h *Handler) postTransfer(c *gin.Context) {
	var input models.TransactionInputHandler
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	data := models.TransactionInput{
		SenderId:    input.SenderId,
		RecipientId: input.RecipientId,
		Amount:      fromNormal(input.Amount),
		Description: input.Description,
	}

	if err := h.services.MakeTransaction(data); err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		senderBalance, err := h.services.Account.GetBalance(data.SenderId)
		if err != nil {
			log.Printf("Error getting account balance during publishing to queue (id=%d)", data.SenderId)
			return
		}

		recipientBalance, err := h.services.Account.GetBalance(data.RecipientId)
		if err != nil {
			log.Printf("Error getting account balance during publishing to queue (id=%d)", data.RecipientId)
			return
		}

		rabbitmq.PublishTransfer(data, senderBalance, recipientBalance)
	}()

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
func (h *Handler) getTransactions(c *gin.Context) {
	var input models.TransactionInputHandlerId

	outputParams, err := Validate(c.Query("limit"), c.Query("offset"), c.Query("sort"))
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := h.services.GetAccountTransactions(input.Id, outputParams)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getTransactionsResponse{
		Data: transactions,
	})
}

func Validate(limit, offset, sort string) (models.OutputParams, error) {
	var limitRes int
	if limit == "" {
		limitRes = 10
	} else {
		var err error
		limitRes, err = strconv.Atoi(limit)
		if err != nil {
			return models.OutputParams{}, fmt.Errorf("Invalid limit param")
		}
	}

	var offsetRes int
	if offset == "" {
		offsetRes = 0
	} else {
		var err error
		offsetRes, err = strconv.Atoi(offset)
		if err != nil {
			return models.OutputParams{}, fmt.Errorf("Invalid offset param")
		}
	}

	var sortCol, sortDir string
	if sort != "" {
		if strings.Contains(sort, ":") {
			sortArr := strings.Split(sort, ":")

			sortCol = sortArr[0]
			sortDir = sortArr[1]

			if (sortCol != "timestamp" && sortCol != "amount") || (sortDir != "asc" && sortDir != "desc") {
				return models.OutputParams{}, fmt.Errorf("Invalid sort param")
			}
		} else if sortCol != "timestamp" && sortCol != "amount" {
			return models.OutputParams{}, fmt.Errorf("Invalid sort param")
		}
	} else {
		sortCol = "timestamp"
		sortDir = "desc"
	}

	return models.OutputParams{
		Limit:   limitRes,
		Offset:  offsetRes,
		SortCol: sortCol,
		SortDir: sortDir,
	}, nil
}
