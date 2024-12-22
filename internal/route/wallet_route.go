package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"go-wallet-service/internal/middleware"
	"go-wallet-service/internal/object"
	"go-wallet-service/internal/service"
	"go-wallet-service/utils"
)

var walletService *service.WalletService

func WalletRoutes(r *mux.Router, us *service.UserService, ws *service.WalletService) {
	walletService = ws

	walletRouter := r.PathPrefix("/wallet").Subrouter()
	walletRouter.Use(middleware.OAuth(us))

	walletRouter.HandleFunc("/{walletId}/balance", balanceHandler).Methods("GET")
	walletRouter.HandleFunc("/{walletId}/deposit", depositHandler).Methods("POST")
	walletRouter.HandleFunc("/{walletId}/withdraw", withdrawHandler).Methods("POST")
}

func balanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestWalletId := vars["walletId"]

	walletId, err := strconv.Atoi(requestWalletId)
	if err != nil {
		log.Error().Str("error", "invalid_parameter").Str("error_description", fmt.Sprintf("%s must be a valid integer", requestWalletId)).Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	wallet, err := walletService.GetById(walletId)
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to fetch wallet [%d]", walletId)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to fetch wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"balance": wallet.Balance,
	})
}

func depositHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestWalletId := vars["walletId"]

	walletId, err := strconv.Atoi(requestWalletId)
	if err != nil {
		log.Error().Str("error", "invalid_parameter").Str("error_description", fmt.Sprintf("%s must be a valid integer", requestWalletId)).Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	requestParams, ok := r.Context().Value("requestParams").(string)
	if !ok {
		log.Error().Str("error", "invalid_parameter").Str("error_description", "Parameters are missing, not expected or not matching the required format").Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	var depositParams object.DepositParams
	err = json.Unmarshal([]byte(requestParams), &depositParams)
	if err != nil {
		log.Error().Str("error", "json_unmarshal_error").Str("error_description", fmt.Sprintf("Unable to unmarshal parameters: %v", requestParams)).Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	wallet, err := walletService.GetById(walletId)
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to fetch wallet [%d]", walletId)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to fetch wallet", http.StatusInternalServerError)
		return
	}

	wallet.Balance += depositParams.Amount

	err = walletService.Update(wallet, depositParams.Amount, "deposit")
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to deposit [%f] to wallet [%d] due to: %v", wallet.Balance, walletId, err.Error())).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to deposit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Deposit Successful",
	})
}

func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestWalletId := vars["walletId"]

	walletId, err := strconv.Atoi(requestWalletId)
	if err != nil {
		log.Error().Str("error", "invalid_parameter").Str("error_description", fmt.Sprintf("%s must be a valid integer", requestWalletId)).Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	requestParams, ok := r.Context().Value("requestParams").(string)
	if !ok {
		log.Error().Str("error", "invalid_parameter").Str("error_description", "Parameters are missing, not expected or not matching the required format").Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	var withdrawParams object.WithdrawParams
	err = json.Unmarshal([]byte(requestParams), &withdrawParams)
	if err != nil {
		log.Error().Str("error", "json_unmarshal_error").Str("error_description", fmt.Sprintf("Unable to unmarshal parameters: %v", requestParams)).Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	wallet, err := walletService.GetById(walletId)
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to fetch wallet [%d]", walletId)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to fetch wallet", http.StatusInternalServerError)
		return
	}

	if wallet.Balance < withdrawParams.Amount {
		utils.HttpErrorResponse(w, "operation_not_permitted", "The withdrawal amount is greater than wallet balance", http.StatusForbidden)
		return
	}

	wallet.Balance -= withdrawParams.Amount

	err = walletService.Update(wallet, withdrawParams.Amount, "withdraw")
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to deposit [%f] to wallet [%d]", wallet.Balance, walletId)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to deposit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Withdraw Successful",
	})
}
