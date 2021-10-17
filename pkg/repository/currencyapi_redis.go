package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var ctx = context.Background()

type CurrencyAPIRedis struct {
	redisClient *redis.Client
}

func NewCurrencyAPIRedis(redisClient *redis.Client) *CurrencyAPIRedis {
	return &CurrencyAPIRedis{redisClient: redisClient}
}

func (c *CurrencyAPIRedis) GetCurrency(code string) (float64, error) {
	value, err := c.redisClient.Get(ctx, code).Result()
	if err == redis.Nil {
		log.Println("Making a request for exchange rates")
		apiKey := viper.GetString("api.api_key")
		baseCurrency := viper.GetString("api.base_currency")
		url := fmt.Sprintf("https://freecurrencyapi.net/api/v2/latest?apikey=%s&base_currency=%s", apiKey, baseCurrency)

		resp, err := http.Get(url)
		if err != nil {
			return 0, err
		}

		data := CurrencyResponse{}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return 0, err
		}
		//сохранение курса валют в кэш
		log.Println("Added exchange rates to cache")
		for key, value := range data.Data {
			err = c.redisClient.Set(ctx, key, value, 3600*time.Second).Err()
			if err != nil {
				return 0, err
			}
		}

		value, err = c.redisClient.Get(ctx, code).Result()
		if err == redis.Nil {
			return 0, err
		}
	}
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

type CurrencyResponse struct {
	Query struct {
		Apikey       string `json:"apikey"`
		BaseCurrency string `json:"base_currency"`
		Timestamp    int    `json:"timestamp"`
	} `json:"query"`
	Data map[string]float64 `json:"data"`
}
