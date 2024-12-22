package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"go-wallet-service/internal/middleware"
	"go-wallet-service/internal/model"
	"go-wallet-service/internal/object"
	"go-wallet-service/internal/service"
	"go-wallet-service/utils"
)

var userService *service.UserService

func UserRoutes(r *mux.Router, us *service.UserService, ws *service.WalletService) {
	userService = us
	walletService = ws

	userRouter := r.PathPrefix("/user/{userId}").Subrouter()
	userRouter.Use(middleware.OAuth(us))

	userRouter.HandleFunc("/transactions", transactionHandler).Methods("GET")
	userRouter.HandleFunc("/transfer", transferHandler).Methods("POST")
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Error().Str("error", "invalid_parameter").Str("error_description", fmt.Sprintf("%s must be a valid integer", userId)).Send()
		return
	}

	transactions, err := userService.GetUserTransactionsByUserId(id)
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to fetch transactions for user [%d]", id)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to fetch transactions", http.StatusInternalServerError)
		return
	}

	var receiverUserIds []int
	for _, txn := range transactions {
		receiverUserIds = append(receiverUserIds, txn.ReceiverUserId)
	}

	users, err := userService.GetByIds(receiverUserIds)
	if err != nil {
		log.Error().Str("error", "user_service_error").Str("error_description", fmt.Sprintf("Unable to fetch users due to: %s ", err.Error())).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to fetch user IDs", http.StatusInternalServerError)
		return
	}

	receiverUsernames := make(map[int]string)
	for _, user := range users {
		receiverUsernames[user.ID] = user.Username
	}

	json.NewEncoder(w).Encode(mapTransactions(transactions, receiverUsernames))
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	requestParams, ok := r.Context().Value("requestParams").(string)
	if !ok {
		log.Error().Str("error", "invalid_parameter").Str("error_description", "Parameters are missing, not expected or not matching the required format").Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	var transferParams object.TransferParams
	err := json.Unmarshal([]byte(requestParams), &transferParams)
	if err != nil {
		log.Error().Str("error", "json_unmarshal_error").Str("error_description", fmt.Sprintf("Unable to unmarshal parameters: %v", requestParams)).Send()
		utils.HttpErrorResponse(w, "invalid_parameter", "Parameters are missing, not expected or not matching the required format", http.StatusBadRequest)
		return
	}

	senderWallet, err := walletService.GetWalletByUserId(transferParams.UserId)
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to fetch wallet for user [%d]", transferParams.UserId)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to fetch wallet", http.StatusInternalServerError)
	}
	if senderWallet[0].Balance < transferParams.Amount {
		utils.HttpErrorResponse(w, "operation_not_permitted", "The transfer amount is greater than wallet balance", http.StatusForbidden)
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Error().Str("error", "invalid_parameter").Str("error_description", fmt.Sprintf("%s must be a valid integer", userId)).Send()
		return
	}

	if id != transferParams.UserId {
		utils.HttpErrorResponse(w, "operation_not_permitted", "The current user is not the sender", http.StatusForbidden)
	}

	receiverWallet, err := walletService.GetWalletByUserId(transferParams.ReceiverUserId)
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to fetch wallet for user [%d]", transferParams.UserId)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to fetch wallet", http.StatusInternalServerError)
	}

	senderWalletBalance := senderWallet[0].Balance - transferParams.Amount
	receiverWalletBalance := receiverWallet[0].Balance + transferParams.Amount

	err = userService.Transfer(transferParams.UserId, senderWallet[0].ID, senderWalletBalance, transferParams.ReceiverUserId, receiverWallet[0].ID, receiverWalletBalance, transferParams.Amount)
	if err != nil {
		log.Error().Str("error", "database_error").Str("error_description", fmt.Sprintf("Unable to transfer from wallet [%d] to wallet [%d]", senderWallet[0].ID, receiverWallet[0].ID)).Send()
		utils.HttpErrorResponse(w, "bad_request", "Unable to deposit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Transfer Successful",
	})
}

func mapTransactions(source []model.Transaction, receivers map[int]string) []model.TransactionData {
	transactionModels := make([]model.TransactionData, len(source))
	for i, txn := range source {
		transactionModels[i] = model.TransactionData{
			TransactionDate: txn.CreationDate,
			Type:            txn.Type,
			Amount:          txn.Amount,
			Receiver:        receivers[txn.ReceiverUserId],
		}
	}
	return transactionModels
}
