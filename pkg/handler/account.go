package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microphoneabuser/balance-service/models"
	"github.com/microphoneabuser/balance-service/rabbitmq"
)

func (h *Handler) getBalance(c *gin.Context) {
	var input models.AccountInputHandler
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	balance, err := h.services.GetBalance(input.Id)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	finalBalance, err := setCurrency(balance, c.Query("currency"), h)
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, balanceResponse{finalBalance})
}
func (h *Handler) postAccrual(c *gin.Context) {
	var input models.AccountInputHandler
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	data := models.AccountInput{
		Id:          input.Id,
		Amount:      fromNormal(input.Amount),
		Description: input.Description,
	}

	if err := h.services.Accrual(data); err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		balance, err := h.services.Account.GetBalance(data.Id)
		if err != nil {
			log.Printf("Error getting account balance during publishing to queue (id=%d)", data.Id)
		} else {
			rabbitmq.PublishAccrualDebiting(data, balance, true)
		}
	}()

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
func (h *Handler) postDebiting(c *gin.Context) {
	var input models.AccountInputHandler
	if err := c.BindJSON(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	data := models.AccountInput{
		Id:          input.Id,
		Amount:      fromNormal(input.Amount),
		Description: input.Description,
	}

	if err := h.services.Debiting(data); err != nil {
		newErrorMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		balance, err := h.services.Account.GetBalance(data.Id)
		if err != nil {
			log.Printf("Error getting account balance during publishing to queue (id=%d)", data.Id)
		} else {
			rabbitmq.PublishAccrualDebiting(data, balance, false)
		}
	}()

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func toNormal(amount int) float64 {
	return float64(amount) / 100
}

func fromNormal(amount float64) int {
	return int(amount * 100)
}

func setCurrency(amount int, currency string, h *Handler) (float64, error) {
	if currency != "" {
		if validCurrency[currency] {
			var err error
			if amount, err = h.services.Currency.Convert(amount, currency); err != nil {
				return 0, err
			}
		} else {
			return 0, fmt.Errorf("Invalid currency param")
		}
	}
	return toNormal(amount), nil
}

var validCurrency = map[string]bool{
	"USD": true, "JPY": true, "CNY": true, "CHF": true, "CAD": true, "MXN": true, "INR": true, "BRL": true, "RUB": true, "KRW": true, "IDR": true, "TRY": true,
	"SAR": true, "SEK": true, "NGN": true, "PLN": true, "ARS": true, "NOK": true, "TWD": true, "IRR": true, "AED": true, "COP": true, "THB": true, "ZAR": true,
	"DKK": true, "MYR": true, "SGD": true, "ILS": true, "HKD": true, "EGP": true, "PHP": true, "CLP": true, "PKR": true, "IQD": true, "DZD": true, "KZT": true,
	"QAR": true, "CZK": true, "PEN": true, "RON": true, "VND": true, "BDT": true, "HUF": true, "UAH": true, "AOA": true, "MAD": true, "OMR": true, "CUC": true,
	"BYR": true, "AZN": true, "LKR": true, "SDG": true, "SYP": true, "MMK": true, "DOP": true, "UZS": true, "KES": true, "GTQ": true, "URY": true, "HRV": true,
	"MOP": true, "ETB": true, "CRC": true, "TZS": true, "TMT": true, "TND": true, "PAB": true, "LBP": true, "RSD": true, "LYD": true, "GHS": true, "YER": true,
	"BOB": true, "BHD": true, "CDF": true, "PYG": true, "UGX": true, "SVC": true, "TTD": true, "AFN": true, "NPR": true, "HNL": true, "BIH": true, "BND": true,
	"ISK": true, "KHR": true, "GEL": true, "MZN": true, "BWP": true, "PGK": true, "JMD": true, "XAF": true, "NAD": true, "ALL": true, "SSP": true, "MUR": true,
	"MNT": true, "NIO": true, "LAK": true, "MKD": true, "AMD": true, "MGA": true, "XPF": true, "TJS": true, "HTG": true, "BSD": true, "MDL": true, "RWF": true,
	"KGS": true, "GNF": true, "SRD": true, "SLL": true, "XOF": true, "MWK": true, "FJD": true, "ERN": true, "SZL": true, "GYD": true, "BIF": true, "KYD": true,
	"MVR": true, "LSL": true, "LRD": true, "CVE": true, "DJF": true, "SCR": true, "SOS": true, "GMD": true, "KMF": true, "STD": true, "XRP": true, "AUD": true,
	"BGN": true, "BTC": true, "JOD": true, "GBP": true, "ETH": true, "EUR": true, "LTC": true, "NZD": true,
}
