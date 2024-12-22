package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go-wallet-service/internal/repository"
	"go-wallet-service/internal/route"
	"go-wallet-service/internal/service"
)

const (
	_DotEnvLogLevel = "LOG_LEVEL"

	_DotEnvDatabaseUrl = "DATABASE_URL"
	_DotEnvBaseUrl     = "BASE_URL"
)

var (
	LogLevel      string
	DatabaseUrl   string
	walletService *service.WalletService
	userService   *service.UserService
)

func main() {
	loadAndCheckDotEnvFile()
	setLogLevel()
	db := connectToDatabase()
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	userService = service.NewUserService(userRepository)
	walletService = service.NewWalletService(walletRepository)

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	route.WalletRoutes(router, userService, walletService)
	route.UserRoutes(router, userService, walletService)

	log.Debug().Msg("Connected to Server")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Panic().Str("error", "connection_failed").Str("error_description", fmt.Sprintf("Unable to connect to server due to: %s", err.Error())).Send()
	}
}

func loadAndCheckDotEnvFile() {
	dotEnvLoaded := false

	configToCheck := []string{
		_DotEnvLogLevel,
		_DotEnvDatabaseUrl,
		_DotEnvBaseUrl,
	}
	for _, config := range configToCheck {
		if os.Getenv(config) == "" {
			err := godotenv.Load()
			if err != nil {
				log.Panic().Str("error", "invalid_dot_env").Str("error_description", ".env file is missing").Send()
			}
			dotEnvLoaded = true
			break
		}
	}
	if dotEnvLoaded {
		for _, config := range configToCheck {
			if os.Getenv(config) == "" {
				log.Panic().Str("error", "invalid_dot_env").Str("error_description", fmt.Sprintf("%s is missing in .env", config)).Send()
			}
		}
	}

	LogLevel = os.Getenv(_DotEnvLogLevel)
}

func connectToDatabase() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic().Str("error", "database_connection_failed").Str("error_description", fmt.Sprintf("Unable to connect to database due to: %s for %v", err.Error(), os.Getenv("DATABASE_URL"))).Send()
	}
	return db
}

func setLogLevel() {
	switch LogLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
}
