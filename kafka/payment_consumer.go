package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type PaymentConsumer struct {
	reader *kafka.Reader
}

func NewPaymentConsumer(brokers []string) *PaymentConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "payments",
		GroupID: "payment-processor",
	})

	return &PaymentConsumer{
		reader: reader,
	}
}

func (c *PaymentConsumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var event map[string]interface{}
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			eventType, ok := event["type"].(string)
			if !ok {
				log.Printf("Invalid event type")
				continue
			}

			paymentData, ok := event["payment"].(map[string]interface{})
			if !ok {
				log.Printf("Invalid payment data")
				continue
			}

			switch eventType {
			case "payment.created":
				// TODO: Process payment based on payment type
				// - For CARD payments: Check card balance and deduct amount
				// - For CASH payments: Check customer balance and deduct amount
				log.Printf("Processing payment: %v", paymentData)

			case "payment.cancelled":
				// Payment is already cancelled, just log it
				log.Printf("Payment cancelled: %v", paymentData)
			}
		}
	}
}

func (c *PaymentConsumer) Close() error {
	return c.reader.Close()
}
