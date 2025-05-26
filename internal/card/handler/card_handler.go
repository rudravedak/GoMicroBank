package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"govo/internal/card/service"
)

type CardHandler struct {
	service *service.CardService
}

func NewCardHandler(service *service.CardService) *CardHandler {
	return &CardHandler{service: service}
}

type CreateCardRequest struct {
	CustomerID  uint    `json:"customer_id"`
	CardNumber  string  `json:"card_number"`
	CardType    string  `json:"card_type"`
	ExpiryDate  string  `json:"expiry_date"`
	CVV         string  `json:"cvv"`
	CreditLimit float64 `json:"credit_limit"`
	Balance     float64 `json:"balance"`
}

type CardResponse struct {
	ID          uint    `json:"id"`
	CustomerID  uint    `json:"customer_id"`
	CardNumber  string  `json:"card_number"`
	CardType    string  `json:"card_type"`
	ExpiryDate  string  `json:"expiry_date"`
	CreditLimit float64 `json:"credit_limit"`
	Balance     float64 `json:"balance"`
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.AddCard(
		req.CustomerID,
		req.CardNumber,
		req.CardType,
		req.ExpiryDate,
		req.CVV,
		req.CreditLimit,
		req.Balance,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	cards, err := h.service.GetCustomerCards(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]CardResponse, len(cards))
	for i, c := range cards {
		response[i] = CardResponse{
			ID:          c.ID,
			CustomerID:  c.CustomerID,
			CardNumber:  c.CardNumber,
			CardType:    c.CardType,
			ExpiryDate:  c.ExpiryDate,
			CreditLimit: c.CreditLimit,
			Balance:     c.Balance,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CardHandler) ListCards(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	cards, err := h.service.GetCustomerCards(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]CardResponse, len(cards))
	for i, c := range cards {
		response[i] = CardResponse{
			ID:          c.ID,
			CustomerID:  c.CustomerID,
			CardNumber:  c.CardNumber,
			CardType:    c.CardType,
			ExpiryDate:  c.ExpiryDate,
			CreditLimit: c.CreditLimit,
			Balance:     c.Balance,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	cardNumber := r.URL.Query().Get("card_number")

	if customerID == "" || cardNumber == "" {
		http.Error(w, "Customer ID and Card Number are required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemoveCard(uint(id), cardNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
