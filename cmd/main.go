package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	balanceservice "github.com/microphoneabuser/balance-service"
	"github.com/microphoneabuser/balance-service/pkg/handler"
	"github.com/microphoneabuser/balance-service/pkg/repository"
	"github.com/microphoneabuser/balance-service/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializtion config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(
		repository.PostgresConfig{
			Host:     viper.GetString("db.Host"),
			Port:     viper.GetString("db.Port"),
			Username: viper.GetString("db.Username"),
			DBName:   viper.GetString("db.DBname"),
			SSLMode:  viper.GetString("db.SSLmode"),
			Password: os.Getenv("DB_PASSWORD"),
		})
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}

	redisClient := repository.NewRedisClient(
		repository.RedisConfig{
			RedisAddr:      viper.GetString("redis.RedisAddr"),
			RedisPassword:  viper.GetString("redis.RedisPassword"),
			RedisDB:        viper.GetString("redis.RedisDB"),
			RedisDefaultdb: viper.GetString("redis.RedisDefaultdb"),
			MinIdleConns:   viper.GetInt("redis.MinIdleConns"),
			PoolSize:       viper.GetInt("redis.PoolSize"),
			PoolTimeout:    viper.GetInt("redis.PoolTimeout"),
			Password:       viper.GetString("redis.Password"),
			DB:             viper.GetInt("redis.DB"),
		})
	err = redisClient.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		log.Fatalf("Failed to initialize redis: %s", err.Error())
	}

	repos := repository.NewRepository(db, redisClient)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(balanceservice.Server)
	if err := srv.Run(viper.GetString("port"), handlers.SetupRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
