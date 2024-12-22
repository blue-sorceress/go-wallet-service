package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWalletBalance(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/wallet/8/balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestBody struct {
			UserId int `json:"userId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if requestBody.UserId == 2 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"balance": 100.0,
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	requestBody := map[string]int{"userId": 2}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	resp, err := http.Post(server.URL+"/wallet/8/balance", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, 100.0, responseBody["balance"])
}

func TestGetUserTransactions(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/user/2/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		transactions := []map[string]interface{}{
			{"id": 1, "amount": 200.00, "type": "deposit"},
			{"id": 2, "amount": 150.00, "type": "withdrawal"},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transactions)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/user/2/transactions")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.Len(t, responseBody, 2)
	assert.Equal(t, 200.00, responseBody[0]["amount"])
	assert.Equal(t, "deposit", responseBody[0]["type"])
	assert.Equal(t, 150.00, responseBody[1]["amount"])
	assert.Equal(t, "withdrawal", responseBody[1]["type"])
}

func TestDepositToWallet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/wallet/8/deposit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestBody struct {
			UserId int     `json:"userId"`
			Amount float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if requestBody.UserId == 2 && requestBody.Amount > 0 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":    "Deposit successful",
				"newBalance": 3000.0,
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	requestBody := map[string]interface{}{
		"userId": 2,
		"amount": 2000,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	resp, err := http.Post(server.URL+"/wallet/8/deposit", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Deposit successful", responseBody["message"])
	assert.Equal(t, 3000.0, responseBody["newBalance"])
}

func TestWithdrawFromWallet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/wallet/8/withdraw", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestBody struct {
			UserId int     `json:"userId"`
			Amount float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if requestBody.UserId == 2 && requestBody.Amount > 0 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":    "Withdrawal successful",
				"newBalance": 1500.00,
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	requestBody := map[string]interface{}{
		"userId": 2,
		"amount": 500.00,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	resp, err := http.Post(server.URL+"/wallet/8/withdraw", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Withdrawal successful", responseBody["message"])
	assert.Equal(t, 1500.00, responseBody["newBalance"])
}

func TestTransferToAnotherUser(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/user/2/transfer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestBody struct {
			UserId         int     `json:"userId"`
			ReceiverUserId int     `json:"receiverUserId"`
			Amount         float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if requestBody.UserId == 2 && requestBody.ReceiverUserId == 1 && requestBody.Amount > 0 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":            "Transfer successful",
				"newSenderBalance":   1000.00, // Example new balance for the sender
				"newReceiverBalance": 570.00,  // Example new balance for the receiver
			})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	requestBody := map[string]interface{}{
		"userId":         2,
		"receiverUserId": 1,
		"amount":         570.00,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	resp, err := http.Post(server.URL+"/user/2/transfer", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Transfer successful", responseBody["message"])
	assert.Equal(t, 1000.00, responseBody["newSenderBalance"])
	assert.Equal(t, 570.00, responseBody["newReceiverBalance"])
}
