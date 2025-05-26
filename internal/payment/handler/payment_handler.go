package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"govo/internal/payment/service"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

type CreatePaymentRequest struct {
	CustomerID  uint    `json:"customer_id"`
	CardID      uint    `json:"card_id"`
	Amount      float64 `json:"amount"`
	PaymentType string  `json:"payment_type"`
	Description string  `json:"description"`
}

type PaymentResponse struct {
	ID          uint      `json:"id"`
	CustomerID  uint      `json:"customer_id"`
	CardID      uint      `json:"card_id"`
	Amount      float64   `json:"amount"`
	PaymentType string    `json:"payment_type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payment, err := h.service.CreatePayment(
		r.Context(),
		req.CustomerID,
		req.CardID,
		req.Amount,
		req.PaymentType,
		req.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PaymentResponse{
		ID:          payment.ID,
		CustomerID:  payment.CustomerID,
		CardID:      payment.CardID,
		Amount:      payment.Amount,
		PaymentType: payment.PaymentType,
		Status:      payment.Status,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get("id")
	if paymentID == "" {
		http.Error(w, "Payment ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(paymentID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := h.service.GetPayment(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PaymentResponse{
		ID:          payment.ID,
		CustomerID:  payment.CustomerID,
		CardID:      payment.CardID,
		Amount:      payment.Amount,
		PaymentType: payment.PaymentType,
		Status:      payment.Status,
		Description: payment.Description,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	status := r.URL.Query().Get("status")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	var customerIDUint uint
	if customerID != "" {
		id, err := strconv.ParseUint(customerID, 10, 32)
		if err == nil {
			customerIDUint = uint(id)
		}
	}

	payments, err := h.service.ListPayments(r.Context(), customerIDUint, status, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		response[i] = PaymentResponse{
			ID:          p.ID,
			CustomerID:  p.CustomerID,
			CardID:      p.CardID,
			Amount:      p.Amount,
			PaymentType: p.PaymentType,
			Status:      p.Status,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PaymentHandler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get("id")
	if paymentID == "" {
		http.Error(w, "Payment ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(paymentID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	reason := r.URL.Query().Get("reason")
	if reason == "" {
		http.Error(w, "Cancel reason is required", http.StatusBadRequest)
		return
	}

	err = h.service.CancelPayment(r.Context(), uint(id), reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
