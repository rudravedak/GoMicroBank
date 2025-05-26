package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"govo/internal/payment/model"
	"govo/internal/payment/repository"
	"govo/kafka"
)

type PaymentService struct {
	repo        *repository.PaymentRepository
	kafkaClient *kafka.Client
}

func NewPaymentService(repo *repository.PaymentRepository, kafkaClient *kafka.Client) *PaymentService {
	return &PaymentService{
		repo:        repo,
		kafkaClient: kafkaClient,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, customerID, cardID uint, amount float64, paymentType, description string) (*model.Payment, error) {
	// Ödeme tipi kontrolü
	if paymentType != "CARD" && paymentType != "CASH" {
		return nil, errors.New("invalid payment type")
	}

	// Kart ödemesi için kart ID kontrolü
	if paymentType == "CARD" && cardID == 0 {
		return nil, errors.New("card ID is required for card payments")
	}

	// Ödeme kaydı oluştur
	payment := &model.Payment{
		CustomerID:  customerID,
		CardID:      cardID,
		Amount:      amount,
		PaymentType: paymentType,
		Status:      "PENDING",
		Description: description,
	}

	// Ödeme kaydını veritabanına kaydet
	payment, err := s.repo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	// Kafka'ya ödeme olayını gönder
	event := map[string]interface{}{
		"event_type":   "PAYMENT_CREATED",
		"payment_id":   payment.ID,
		"customer_id":  payment.CustomerID,
		"card_id":      payment.CardID,
		"amount":       payment.Amount,
		"payment_type": payment.PaymentType,
		"status":       payment.Status,
		"description":  payment.Description,
		"created_at":   payment.CreatedAt,
	}

	if err := s.kafkaClient.SendMessage("payments", event); err != nil {
		// Kafka hatası ödemeyi etkilemesin, sadece logla
		fmt.Printf("Failed to send payment event to Kafka: %v\n", err)
	}

	return payment, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, id uint) (*model.Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PaymentService) ListPayments(ctx context.Context, customerID uint, status string, startDate, endDate *time.Time) ([]*model.Payment, error) {
	return s.repo.List(ctx, customerID, status, startDate, endDate)
}

func (s *PaymentService) CancelPayment(ctx context.Context, id uint, reason string) error {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	// Sadece PENDING veya PROCESSING durumundaki ödemeler iptal edilebilir
	if payment.Status != "PENDING" && payment.Status != "PROCESSING" {
		return errors.New("only pending or processing payments can be cancelled")
	}

	// Ödeme durumunu güncelle
	payment.Status = "CANCELLED"
	payment.Description = fmt.Sprintf("Cancelled: %s", reason)

	if err := s.repo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to cancel payment: %v", err)
	}

	// Kafka'ya iptal olayını gönder
	event := map[string]interface{}{
		"event_type":   "PAYMENT_CANCELLED",
		"payment_id":   payment.ID,
		"customer_id":  payment.CustomerID,
		"card_id":      payment.CardID,
		"amount":       payment.Amount,
		"payment_type": payment.PaymentType,
		"status":       payment.Status,
		"description":  payment.Description,
		"cancelled_at": time.Now(),
	}

	if err := s.kafkaClient.SendMessage("payments", event); err != nil {
		// Kafka hatası işlemi etkilemesin, sadece logla
		fmt.Printf("Failed to send cancellation event to Kafka: %v\n", err)
	}

	return nil
}
