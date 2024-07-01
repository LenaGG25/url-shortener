package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"url-shortener/client/internal/config"
	"url-shortener/client/internal/handlers"
	"url-shortener/client/internal/service"
	"url-shortener/internal/pkg/pb"
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
		ShortURLPrefix: viper.GetString("service.short_url_prefix"),
		Port:           viper.GetInt("service.port"),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	grpcClient, err := grpc.NewClient(
		os.Getenv("URL_SHORTENER_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect with URL Shortener grpc service: %v", err)
	}
	defer grpcClient.Close()

	srv := service.New(
		pb.NewURLShortenerClient(grpcClient),
		fmt.Sprintf("%s:%d", cfg.ShortURLPrefix, cfg.Port),
	)
	handler := handlers.New(srv)

	http.Handle("/", createRouter(handler))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	go func() {
		log.Println("Server is running on port:", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Gracefully shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server shutdown is successful!")
}

func createRouter(handler *handlers.Handler) *mux.Router {
	router := mux.NewRouter()
	router.Use(logMiddleware)

	router.HandleFunc("/shorten", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			handler.CreateURL(w, req)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	router.HandleFunc("/short_url/{short_url}", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			handler.GetOriginalByShortURL(w, req)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	router.HandleFunc("/stats/{short_url}", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			handler.GetStatsByShortURL(w, req)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	return router
}

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var buf bytes.Buffer
		tee := io.TeeReader(request.Body, &buf)
		body, err := io.ReadAll(tee)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		request.Body = io.NopCloser(&buf)

		log.Printf(
			"URL: %s\nMethod: %s\nQueryParams: %v\nBody: %v",
			request.URL.String(),
			request.Method,
			mux.Vars(request),
			body,
		)

		handler.ServeHTTP(writer, request)
	})
}

func initConfig() error {
	viper.AddConfigPath(configsPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
