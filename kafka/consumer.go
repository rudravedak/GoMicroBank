package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	topics   []string
}

func NewConsumer(brokers []string) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Retry/backoff mekanizması
	maxRetries := 30
	retryInterval := 2 * time.Second

	var consumer sarama.Consumer
	var err error

	for i := 0; i < maxRetries; i++ {
		consumer, err = sarama.NewConsumer(brokers, config)
		if err == nil {
			break
		}

		log.Printf("Kafka consumer oluşturulamadı (deneme %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("%v sonra tekrar denenecek...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		panic(fmt.Sprintf("Kafka consumer oluşturulamadı (%d deneme sonrası): %v", maxRetries, err))
	}

	log.Println("Kafka consumer başarıyla oluşturuldu!")
	return &Consumer{
		consumer: consumer,
		topics:   []string{"payments"},
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for _, topic := range c.topics {
		partitions, err := c.consumer.Partitions(topic)
		if err != nil {
			log.Printf("Failed to get partitions for topic %s: %v", topic, err)
			continue
		}

		for _, partition := range partitions {
			pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Failed to start consumer for partition %d: %v", partition, err)
				continue
			}

			wg.Add(1)
			go func(pc sarama.PartitionConsumer) {
				defer wg.Done()
				defer pc.Close()

				for {
					select {
					case msg := <-pc.Messages():
						var event map[string]interface{}
						if err := json.Unmarshal(msg.Value, &event); err != nil {
							log.Printf("Failed to unmarshal message: %v", err)
							continue
						}

						// Event tipine göre işlem yap
						switch event["event_type"] {
						case "PAYMENT_CREATED":
							handlePaymentCreated(event)
						case "PAYMENT_CANCELLED":
							handlePaymentCancelled(event)
						default:
							log.Printf("Unknown event type: %s", event["event_type"])
						}

					case err := <-pc.Errors():
						log.Printf("Error: %v", err)

					case <-ctx.Done():
						return
					}
				}
			}(pc)
		}
	}

	wg.Wait()
}

func handlePaymentCreated(event map[string]interface{}) {
	// Ödeme oluşturulduğunda yapılacak işlemler
	paymentType := event["payment_type"].(string)
	amount := event["amount"].(float64)
	customerID := uint(event["customer_id"].(float64))

	if paymentType == "CARD" {
		// Kart bakiyesini güncelle
		cardID := uint(event["card_id"].(float64))
		updateCardBalance(cardID, amount)
	} else {
		// Kişisel bakiyeyi güncelle
		updateCustomerBalance(customerID, amount)
	}
}

func handlePaymentCancelled(event map[string]interface{}) {
	// Ödeme iptal edildiğinde yapılacak işlemler
	paymentType := event["payment_type"].(string)
	amount := event["amount"].(float64)
	customerID := uint(event["customer_id"].(float64))

	if paymentType == "CARD" {
		// Kart bakiyesini geri al
		cardID := uint(event["card_id"].(float64))
		refundCardBalance(cardID, amount)
	} else {
		// Kişisel bakiyeyi geri al
		refundCustomerBalance(customerID, amount)
	}
}

// TODO: Bu fonksiyonlar gRPC çağrıları ile implement edilecek
func updateCardBalance(cardID uint, amount float64) {
	// Kart bakiyesini güncelle
	log.Printf("Updating card balance for card %d: -%.2f", cardID, amount)
}

func updateCustomerBalance(customerID uint, amount float64) {
	// Kişisel bakiyeyi güncelle
	log.Printf("Updating customer balance for customer %d: -%.2f", customerID, amount)
}

func refundCardBalance(cardID uint, amount float64) {
	// Kart bakiyesini geri al
	log.Printf("Refunding card balance for card %d: +%.2f", cardID, amount)
}

func refundCustomerBalance(customerID uint, amount float64) {
	// Kişisel bakiyeyi geri al
	log.Printf("Refunding customer balance for customer %d: +%.2f", customerID, amount)
}
