package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"url-shortener/internal/config"
	urlgrpc "url-shortener/internal/grpc"
	"url-shortener/internal/pkg/pb"
	service2 "url-shortener/internal/service"
	"url-shortener/internal/storage/db"
	postgres2 "url-shortener/internal/storage/repository/postgres"
	rediscache "url-shortener/internal/storage/repository/redis"
)

const configsPath = "configs"

func main() {
	if err := godotenv.Load(fmt.Sprintf("%s/.env", configsPath)); err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}

	if err := initConfig(); err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	cfg := &config.Config{
		MaxShortKeySize: viper.GetInt("shortener.max_short_key_size"),
		Port:            viper.GetInt("grpc.port"),
		PostgresHost:    viper.GetString("db.postgres_host"),
		RedisHost:       viper.GetString("redis.host"),
		Expiration:      viper.GetString("redis.expiration"),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer func(lis net.Listener) {
		_ = lis.Close()
	}(lis)

	cacheExpiration, err := time.ParseDuration(cfg.Expiration)
	if err != nil {
		log.Fatalf("Failed to parse duration: %v", err)
	}

	cache := rediscache.NewRedis(
		cacheExpiration,
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, os.Getenv("REDIS_PORT")),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
	)

	database, transactionManager, err := db.NewDB(ctx, generateDsn(cfg))
	if err != nil {
		log.Fatalf("Failed connect to db: %v", err)
	}
	defer transactionManager.GetQueryEngine(ctx).Close()

	repoURL := postgres2.NewURLRepo(database, transactionManager)
	repoURLStats := postgres2.NewURLStatsRepo(database, transactionManager)

	service := service2.New(
		repoURL,
		repoURLStats,
		transactionManager,
		cache,
		service2.NewURLGenerator(cfg.MaxShortKeySize),
	)

	srv := urlgrpc.NewURLServer(service)

	grpcServer := grpc.NewServer()
	pb.RegisterURLShortenerServer(grpcServer, srv)

	go func() {
		log.Println("Server is running on port:", cfg.Port)
		log.Fatal(grpcServer.Serve(lis))
	}()

	<-ctx.Done()

	log.Println("Gracefully shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server shutdown is successful!")
}

func initConfig() error {
	viper.AddConfigPath(configsPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}

func generateDsn(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
}
