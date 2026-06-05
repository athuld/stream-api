package datasource

import (
	"context"
	"database/sql"
	"os"
	"streamapi/utils/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	Client      *sql.DB
	RedisClient *redis.Client
)

func init() {

	godotenv.Load()
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dbIP := os.Getenv("DB_IP")
	logger.Info.Println("=== Setting up database ===")
	var err error
	Client, err = sql.Open("mysql", username+":"+password+"@tcp("+dbIP+")/streamdb")

	if err != nil {
		logger.Error.Println("=== Databse Setup error ===")
	}

	logger.Info.Println("=== Database Successfully Connected ===")
	logger.Info.Println("=== Setting up Redis Client ===")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.Debug.Fatal("=== Redis setup error ===")
	}
	logger.Info.Println("=== Redis client ready ===")
}

func getRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // fallback for local dev without docker
	}
	return addr
}
